import { useState } from "react";

export function Order() {
	const [orderId, setOrderId] = useState("");
	const [orderDetails, setOrderDetails] = useState(null);
	const [error, setError] = useState(null);
	const [loading, setLoading] = useState(false);

	const handleSearch = async () => {
		if (!orderId || isNaN(Number(orderId))) {
			setError("Неверный ID заказа");
			setOrderDetails(null);
			return;
		}

		setLoading(true);
		setError(null);
		setOrderDetails(null);

		try {
			const response = await fetch(
				`http://localhost:8080/api/orders/${orderId}`,
				{
					method: "GET",
					headers: {
						"Content-Type": "application/json",
					},
					credentials: "include",
				}
			);

			const data = await response.json();

			if (response.ok) {
				setOrderDetails(data);
			} else {
				setError(data.error || "Заказ с таким ID не найден");
			}
		} catch {
			setError("Произошла ошибка при поиске заказа.");
		} finally {
			setLoading(false);
		}
	};

	return (
		<>
			<div className="header">
				<h1>Посмотреть детали заказа</h1>
			</div>
			<div className="search">
				<input
					type="text"
					placeholder="Введите ID заказа..."
					value={orderId}
					onChange={(e) => setOrderId(e.target.value)}
				/>
				<button type="submit" onClick={handleSearch} disabled={loading}>
					{loading ? "Поиск..." : "Поиск"}
				</button>
			</div>
			<div className="search-results">
				{error && (
					<div>
						<h2>Упс...</h2>
						<p>{error}</p>
					</div>
				)}
				{orderDetails && (
					<div>
						<h2>Детали заказа</h2>
						<p>Заказ {orderId}</p>
						<pre>
							<code>{JSON.stringify(orderDetails, null, 2)}</code>
						</pre>
					</div>
				)}
			</div>
		</>
	);
}
