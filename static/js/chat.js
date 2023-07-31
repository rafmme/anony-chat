try {
  const token = localStorage.getItem("userToken");
  const { sub } = jwt_decode(token);

  const userTag = document.getElementById("userTag");
  const logoutBtn = document.getElementById("logoutBtn");
  const sendBtn = document.getElementById("sendBtn");
  const onlineUsersHolder = document.getElementById("online-peeps");
  const chatHolder = document.getElementById("chat-holder");

  const userId = `user-${sub.substring(0, 8)}`;
  userTag.innerText = `Your User ID: ${userId}`;

  // Replace 'your_websocket_endpoint' with the actual WebSocket URL
  const websocketUrl = `ws://localhost:6777/ws/chat?auth_token=${encodeURI(
    token
  )}`;

  // Connect to the WebSocket endpoint with the authorization header
  const socket = new WebSocket(websocketUrl);

  // Handle incoming messages
  socket.onmessage = function (event) {
    const data = JSON.parse(event.data);
    const onlineUsers = Number.parseInt(data.clientCount, 10) - 1 || 0;

    if (
      data.msgType.toLowerCase() === "count" &&
      data.clientCount !== undefined
    ) {
      document.getElementById("clientCount").innerText = onlineUsers;

      if (onlineUsers < 1) {
        onlineUsersHolder.innerHTML = "";
      }
    }

    if (
      data.msgType.toLowerCase() === "count" &&
      data.clientsList !== undefined &&
      onlineUsers > 0
    ) {
      for (const user of Object.keys(data.clientsList)) {
        if (user !== userId) {
          onlineUsersHolder.innerHTML = "";
          onlineUsersHolder.innerHTML += `<a href="#" class="list-group-item list-group-item-action border-0">
            <div class="d-flex align-items-start">
              <img
                src="https://bootdey.com/img/Content/avatar/avatar3.png"
                class="rounded-circle mr-1 us"
                alt=${user}
                width="40"
                height="40"
              />
            <div class="flex-grow-1 ml-3">
              ${user}
              <div class="small">
                <span class="fas fa-circle chat-online"></span> Online
              </div>
            </div>
           </div>
          </a>`;
        }
      }
    }

    if (data.msgType.toLowerCase() === "msg") {
      const { message, sender, date: msgTime, mentioned, private } = data;
      const date = new Date(msgTime).toLocaleString();

      if (sender === userId) {
        chatHolder.innerHTML += `<div class="chat-message-right pb-4">
                    <div>
                      <div class="small text-nowrap mt-2 chat-online">
                        ${private ? "private message" : ""}
                      </div>
                      <div class="text-muted small text-nowrap mt-2">
                      ${date}
                    </div>
                    </div>
                    <div class="flex-shrink-1 bg-light rounded py-2 px-3 mr-3">
                      <div class="font-weight-bold mb-1">You</div>
                      ${message}
                    </div>
                  </div>`;
      } else {
        chatHolder.innerHTML += `<div class="chat-message-left pb-4">
        <div>
          <div class="small text-nowrap mt-2 chat-online">
          ${private ? "private message" : ""}
          </div>
          <div class="small text-nowrap mt-2 chat-offline">
          ${mentioned ? "You were mentioned here!" : ""}
          </div>
          <div class="text-muted small text-nowrap mt-2">
          ${date}
        </div>
        </div>
        <div class="flex-shrink-1 bg-light rounded py-2 px-3 ml-3">
          <div class="font-weight-bold mb-1">${sender}</div>
          ${message}
        </div>
      </div>`;
      }
    }
  };

  // Handle connection open event
  socket.onopen = function (event) {
    console.log("WebSocket connection opened.", event);
  };

  // Handle connection error event
  socket.onerror = function (event) {
    console.error("WebSocket error:", event);
  };

  // Handle connection close event
  socket.onclose = function (event) {
    console.log("WebSocket connection closed:", event);
  };

  sendBtn.addEventListener("click", (evt) => {
    evt.preventDefault();
    const textElement = document.getElementById("text-msg");
    const textMessage = textElement.value;
    chatMessage = textMessage.trim();

    if (chatMessage !== "" && chatMessage !== " ") {
      socket.send(chatMessage);
      textElement.value = "";
    }
  });

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
