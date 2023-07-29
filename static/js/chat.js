try {
  const token = localStorage.getItem("userToken");
  const { sub } = jwt_decode(token);

  const userTag = document.getElementById("userTag");
  const logoutBtn = document.getElementById("logoutBtn");

  userTag.innerText = `Current User ID: ${sub.substring(0, 8)}`;

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
