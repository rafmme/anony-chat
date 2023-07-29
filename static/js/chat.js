try {
  const token = localStorage.getItem("userToken");
  const { sub } = jwt_decode(token);

  const userTag = document.getElementById("userTag");
  const logoutBtn = document.getElementById("logoutBtn");

  const userId = sub.substring(0, 8)
  userTag.innerText = `Current User ID: ${userId}`;

  // Replace 'your_websocket_endpoint' with the actual WebSocket URL
  const websocketUrl = `ws://localhost:6777/ws/chat?auth_token=${encodeURI(
    token
  )}`;

  // Connect to the WebSocket endpoint with the authorization header
  const socket = new WebSocket(websocketUrl);

  // Handle incoming messages
  socket.onmessage = function (event) {
    const messagesDiv = document.getElementById("messages");
    messagesDiv.innerHTML += `<p>${event.data}</p>`;
    console.log(event);
  };

  // Handle connection open event
  socket.onopen = function (event) {
    console.log("WebSocket connection opened.", event);
    socket.send(`Hi, I am user ${userId}`);
  };

  // Handle connection error event
  socket.onerror = function (event) {
    console.error("WebSocket error:", event);
  };

  // Handle connection close event
  socket.onclose = function (event) {
    console.log("WebSocket connection closed:", event);
  };

  const logout = async (evt) => {
    evt.preventDefault();
    localStorage.removeItem("userToken");
    window.location.reload();
  };

  logoutBtn.addEventListener("click", logout);
} catch (error) {
  localStorage.removeItem("userToken");
  window.location.href = "index.html";
}

window.onload = () => {
  try {
    const token = localStorage.getItem("userToken");
    const { sub } = jwt_decode(token);

    if (!sub) {
      localStorage.removeItem("userToken");
      window.location.reload();
    }
  } catch (error) {
    localStorage.removeItem("userToken");
    window.location.href = "index.html";
  }
};
