import os
from locust import HttpUser, task, between

PRODUCT_SERVICE_HOST = os.getenv("PRODUCT_SERVICE_URL", "http://product-service.apps.svc.cluster.local:8001")
CART_SERVICE_HOST = os.getenv("CART_SERVICE_URL", "http://cart-order-service.apps.svc.cluster.local:8002")
PAYMENT_SERVICE_HOST = os.getenv("PAYMENT_SERVICE_URL", "http://payment-service.apps.svc.cluster.local:8003")

class WebsiteUser(HttpUser):
    wait_time = between(1, 5)

    def on_start(self):
        self.cart_id = None
        self.user_id = "user_" + os.urandom(4).hex()

    @task(3)
    def browse_products(self):
        # Visit product listing
        self.client.get(f"{PRODUCT_SERVICE_HOST}/api/products", name="/api/products")

    @task(1)
    def checkout_flow(self):
        # 1. Create Cart
        with self.client.post(f"{CART_SERVICE_HOST}/api/carts", json={"userId": self.user_id}, name="/api/carts [Create]", catch_response=True) as response:
            if response.status_code == 200:
                self.cart_id = response.json().get("id")
            else:
                response.failure(f"Failed to create cart: {response.text}")
                return

        # 2. Add Item to Cart (Assuming product ID '101' exists or generic)
        # Realistically we should pick a product from browse_products, but for simplicity:
        item = {
            "productId": "101",
            "productName": "Test Product",
            "price": 50.0,
            "quantity": 1,
            "selectedSize": "M",
            "selectedColor": "Blue"
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
            "amount": 50.0,
            "currency": "USD",
            "cardNumber": "4242424242424242",
            "cvv": "123",
            "expiryDate": "12/25"
        }
        self.client.post(f"{PAYMENT_SERVICE_HOST}/api/payments", json=payment_payload, name="/api/payments [Pay]")
