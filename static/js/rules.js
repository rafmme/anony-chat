try {
  const proceedBtn = document.getElementById("proceedBtn");
  const backBtn = document.getElementById("backBtn");

  const backAction = async (evt) => {
    evt.preventDefault();
    localStorage.setItem("accepted", false);
    localStorage.removeItem("userToken");
    window.location.reload();
  };

  const proceedAction = async (evt) => {
    evt.preventDefault();
    localStorage.setItem("accepted", true);
    window.location.href = "chat.html";
  };

  backBtn.addEventListener("click", backAction);
  proceedBtn.addEventListener("click", proceedAction);
} catch (error) {
  localStorage.setItem("accepted", false);
  localStorage.removeItem("userToken");
  window.location.href = "index.html";
}

window.onload = () => {
  try {
    const token = localStorage.getItem("userToken");
    const { sub } = jwt_decode(token);

    if (sub && localStorage.getItem("accepted") === true) {
      window.location.href = "chat.html";
    }

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
