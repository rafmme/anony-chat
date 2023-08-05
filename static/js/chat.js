try {
  const token = localStorage.getItem("userToken");
  const { sub } = jwt_decode(token);

  const modal = document.getElementById("myModal");
  const userTag = document.getElementById("userTag");
  const logoutBtn = document.getElementById("logoutBtn");
  const sendBtn = document.getElementById("sendBtn");
  const onlineUsersHolder = document.getElementById("online-peeps");
  const chatHolder = document.getElementById("chat-holder");
  const notificationBell = document.getElementById("notice-bell");
  const mentionText = document.getElementById("m-txt");

  const userId = `user-${sub.substring(0, 8)}`;
  userTag.innerText = `Your User ID: ${userId}`;

  const highlightMentionedId = (message) => {
    if (message.startsWith("@")) {
      const mentionedUserId = message.split(" ")[0];

      return message.replace(
        mentionedUserId,
        `<span class="highlight">${mentionedUserId}</span>`
      );
    }

    return message;
  };

  const websocketUrl = `ws://${window.location.origin
    .replace("https://", "")
    .replace("http://", "")}/ws/chat?auth_token=${encodeURI(token)}`;

  const socket = new WebSocket(websocketUrl);

  // Handle incoming messages
  socket.onmessage = function (event) {
    const data = JSON.parse(event.data);
    const onlineUsers = Number.parseInt(data.clientCount, 10) - 1 || 0;

    if (data.msgType.toLowerCase() === "count") {
      document.getElementById("clientCount").innerText = onlineUsers;

      onlineUsersHolder.innerHTML = "";
      for (const client of Object.keys(data.clientsList)) {
        if (client !== userId) {
          onlineUsersHolder.innerHTML += `<a href="#" class="list-group-item list-group-item-action border-0">
          <div class="d-flex align-items-start">
            <img
              src="https://bootdey.com/img/Content/avatar/avatar3.png"
              class="rounded-circle mr-1 us"
              alt=${client}
              width="40"
              height="40"
            />
          <div class="flex-grow-1 ml-3">
            ${client}
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
      const {
        message,
        action,
        sender,
        date: msgTime,
        mentioned,
        private,
        to
      } = data;
      const date = new Date(msgTime).toLocaleString();

      if (mentioned) {
        notificationBell.classList.remove("hide");
        notificationBell.classList.replace(
          "btn-outline-secondary",
          "bg-danger"
        );
        notificationBell.classList.add("text-white");
        mentionText.innerText = `${sender} just mentioned you.`;

        setTimeout(() => {
          notificationBell.classList.remove("text-white");
          notificationBell.classList.replace(
            "bg-danger",
            "btn-outline-secondary"
          );
          notificationBell.classList.add("hide");
        }, 10000);
      }

      if (sender === userId) {
        chatHolder.innerHTML += `<div class="chat-message-right pb-4">
                    <div>
                      <div class="text-muted small text-nowrap mt-2">
                      ${date}
                    </div>
                      <div class="small text-nowrap mt-2 chat-online">
                        ${private ? `private message to ${to}` : ""}
                      </div>
                    </div>
                    <div class="flex-shrink-1 you-chat-bg rounded py-2 px-3 mr-3">
                      <div class="font-weight-bold mb-1 you-name-bg">You</div>
                      ${highlightMentionedId(message)}
                    </div>
                  </div>`;
      } else {
        chatHolder.innerHTML += `<div class="chat-message-left pb-4">
          <div>
            <div class="text-muted small text-nowrap mt-2">
            ${date}
          </div>
          <div class="small text-nowrap mt-2 chat-online">
          ${private ? `private message from ${sender}` : ""}
          </div>
          <div class="small text-nowrap mt-2 chat-offline">
          ${mentioned ? "You were mentioned here!" : ""}
          </div>
          </div>
          <div class="flex-shrink-1 them-chat-bg rounded py-2 px-3 ml-3">
            <div class="font-weight-bold mb-1 them-name-bg">${sender}</div>
            ${
              message.includes("joined") && message.includes(`@${userId}`)
                ? `You have ${action} the chat`
                : highlightMentionedId(message)
            }
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
    localStorage.setItem("accepted", false);
    localStorage.removeItem("userToken");
    window.location.href = "index.html";
  };

  // Handle connection close event
  socket.onclose = function (event) {
    console.log("WebSocket connection closed:", event);
    localStorage.setItem("accepted", false);
    localStorage.removeItem("userToken");
    window.location.href = "index.html";
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

  function closeModal() {
    modal.style.display = "none";
  }

  function performAction() {
    closeModal();

    localStorage.setItem("accepted", false);
    localStorage.removeItem("userToken");
    window.location.reload();
  }

  const logout = async (evt) => {
    evt.preventDefault();
    modal.style.display = "block";
  };

  logoutBtn.addEventListener("click", logout);
} catch (error) {
  localStorage.setItem("accepted", false);
  localStorage.removeItem("userToken");
  window.location.href = "index.html";
}

window.onload = () => {
  try {
    const token = localStorage.getItem("userToken");
    const { sub } = jwt_decode(token);

    if (!sub || localStorage.getItem("accepted") === false) {
      localStorage.removeItem("userToken");
      window.location.reload();
    }
  } catch (error) {
    localStorage.setItem("accepted", false);
    localStorage.removeItem("userToken");
    window.location.href = "index.html";
  }
};
