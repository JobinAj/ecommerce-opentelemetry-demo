import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Home } from './components/Home';
import { VersaceLogin } from './components/VersaceLogin';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<VersaceLogin />} />
        <Route path="/home" element={<Home />} />
      </Routes>
    </Router>
  );
}

export default App;
