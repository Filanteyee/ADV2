document.getElementById("loadOrdersBtn").addEventListener("click", loadOrders);

function loadOrders() {
    fetch("/orders")
        .then(res => res.json())
        .then(orders => {
            const container = document.getElementById("ordersContainer");
            container.innerHTML = "";

            orders.forEach(order => {
                const card = document.createElement("div");
                card.className = "order-card";
                card.addEventListener("click", () => showOrderDetails(order.id));

                card.innerHTML = `
          <h3>Заказ #${order.id}</h3>
          <p>Статус: <strong>${order.status}</strong></p>
          <div class="order-items">
            ${order.items.map(item => `
              <div class="order-item">
                <img src="${item.image_url}" width="100" /><br>
                ${item.product_name}<br>
                Кол-во: ${item.quantity}
              </div>
            `).join("")}
          </div>
          <div class="order-actions">
            <button onclick="event.stopPropagation(); changeStatus('${order.id}')">Изменить статус</button>
            <button class="cancel-btn" onclick="event.stopPropagation(); cancelOrder('${order.id}')">Отменить заказ</button>
          </div>
        `;
                container.appendChild(card);
            });
        });
}

function changeStatus(orderId) {
    const newStatus = prompt("Введите новый статус (pending, completed, cancelled):");
    if (!newStatus) return;

    fetch(`/orders/${orderId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: newStatus })
    }).then(() => loadOrders());
}

function cancelOrder(orderId) {
    if (!confirm("Вы уверены, что хотите отменить заказ?")) return;

    fetch(`/orders/${orderId}`, {
        method: "DELETE"
    }).then(() => loadOrders());
}

function showOrderDetails(orderId) {
    window.location.href = `/orders/${orderId}/view`;
}
