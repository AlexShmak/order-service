import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Register } from "./pages/register";
import { Login } from "./pages/login";
import { Order } from "./pages/order";

function App() {
	return (
		<Router>
			<Routes>
				<Route path="/" element={<Register />} />
				<Route path="/login" element={<Login />} />
				<Route path="/order" element={<Order />} />
			</Routes>
		</Router>
	);
}

export default App;
