import os
from locust import HttpUser, task, between

PRODUCT_SERVICE_HOST = os.getenv("PRODUCT_SERVICE_URL", "http://product-service.apps.svc.cluster.local:8001")
CART_SERVICE_HOST = os.getenv("CART_SERVICE_URL", "http://cart-order-service.apps.svc.cluster.local:8002")
PAYMENT_SERVICE_HOST = os.getenv("PAYMENT_SERVICE_URL", "http://payment-service.apps.svc.cluster.local:8003")

class WebsiteUser(HttpUser):
    wait_time = between(1, 5)

    def on_start(self):
        self.cart_id = None
        # Generate random user credentials
        random_hex = os.urandom(4).hex()
        self.user_email = f"user_{random_hex}@example.com"
        self.user_password = "password123"
        self.user_name = f"User {random_hex}"
        
        # Register the user
        signup_payload = {
            "name": self.user_name,
            "email": self.user_email,
            "password": self.user_password
        }
        # We need to register against the AUTH service (which is cart-order-service in this setup)
        # Using CART_SERVICE_HOST as it maps to the same service
        with self.client.post(f"{CART_SERVICE_HOST}/api/signup", json=signup_payload, name="/api/signup", catch_response=True) as response:
            if response.status_code == 200:
                self.user = response.json()
                self.user_id = self.user.get("id")
            elif response.status_code == 409:
                # User exists (unlikely with random hex but handled)
                pass 
            else:
                print(f"Failed to signup user: {response.text}")
                self.user_id = "guest" # Fallback, might fail if guest not supported

    @task(3)
    def browse_products(self):
        # Visit product listing
        self.client.get(f"{PRODUCT_SERVICE_HOST}/api/products", name="/api/products")

    @task(1)
    def checkout_flow(self):
        if not self.user_id or self.user_id == "guest":
             # Try to register again if failed previously or skip
             return

        # 1. Create Cart
        with self.client.post(f"{CART_SERVICE_HOST}/api/carts", json={"userId": self.user_id}, name="/api/carts [Create]", catch_response=True) as response:
            if response.status_code == 200:
                self.cart_id = response.json().get("id")
            else:
                response.failure(f"Failed to create cart: {response.text}")
                return

        # 2. Add Item to Cart (Using valid product ID '1')
        item = {
            "productId": "1",
            "productName": "Silk Baroque Shirt",
            "price": 1200.0,
            "quantity": 1,
            "selectedSize": "M",
            "selectedColor": "Black"
        }
        self.client.post(f"{CART_SERVICE_HOST}/api/carts/{self.cart_id}/items", json=item, name="/api/carts/{id}/items [Add]")

        # 3. Create Order
        order_id = None
        with self.client.post(f"{CART_SERVICE_HOST}/api/carts/{self.cart_id}/orders", name="/api/carts/{id}/orders [Create]", catch_response=True) as response:
            if response.status_code == 200:
                order_id = response.json().get("id")
            else:
                response.failure(f"Failed to create order: {response.text}")
                return

        # 4. Process Payment
        payment_payload = {
            "orderId": order_id,
            "amount": 1200.0,
            "currency": "USD",
            "cardNumber": "4242424242424242",
            "cvv": "123",
            "expiryDate": "12/25"
        }
        self.client.post(f"{PAYMENT_SERVICE_HOST}/api/payments", json=payment_payload, name="/api/payments [Pay]")
