const loginForm = document.getElementById("loginForm");
const errorDiv = document.getElementById("loginErrorDiv");
const errorMessage = document.getElementById("loginErrorMessage");

const signupForm = document.getElementById("signupForm");
const signupErrorDiv = document.getElementById("signupErrorDiv");
const signupErrorMessage = document.getElementById("signupErrorMessage");

/**
 * @description a function to validate sign in data
 * @param {String} password
 * @returns {Object} an objects containing the error messages
 */
const validateLoginInput = (password) => {
  const errors = {};

  if (password.trim() === "") {
    errors.password = "Password field is empty";
  }

  const checkPassword =
    password.trim() !== "" && (password.length < 8 || password.length > 26);

  if (checkPassword) {
    errors.password =
      "Password should be within the range of 8 - 26 characters";
  }

  return errors;
};

/**
 * @description a function to validate sign in data
 * @param {String} password
 * @returns {Object} an objects containing the error messages
 */
const validateSignupInput = (password, confirmPassword) => {
  const errors = {};

  if (password.trim() === "") {
    errors.password = "Password field is empty";
  }

  if (password !== confirmPassword) {
    errors.confirmPassword =
      "Confirm Password field does not match Password field input";
  }

  const checkPassword =
    password.trim() !== "" && (password.length < 8 || password.length > 26);

  if (checkPassword) {
    errors.password =
      "Password should be within the range of 8 - 26 characters";
  }

  return errors;
};

/**
 * @description a function to handle the display of error messages
 * @param {Object} errorsObject an Object containing errors
 * @param {Object} divHolder a div tag that houses the error P tag
 * @param {Object} pHolder a P tag that contains the error message
 * @returns {undefined}
 */
const showErrorMessages = (errorsObject, divHolder, pHolder) => {
  let err = '<b><i class="fa fa-warning"></i> Error!!!</b><br>';

  if (errorsObject.length > 0) {
    errorsObject.forEach((errObj) => {
      Object.values(errObj).forEach((e) => {
        err = `${err + e};<br>`;
      });
    });
  } else {
    Object.values(errorsObject).forEach((e) => {
      err = `${err + e};<br>`;
    });
  }

  pHolder.innerHTML = err;
  divHolder.style.opacity = "1";
  divHolder.style.display = "block";

  setTimeout(() => {
    divHolder.style.display = "none";
  }, 5000);
};

/**
 * @description a function to handle the logging in of user into their account on the app
 * @param {Object} evt
 * @returns {undefined}
 */
const login = async (evt) => {
  evt.preventDefault();
  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;
  const loginBtn = document.getElementById("login-btn");
  /* loginBtn.style.background = 'darkgray'; */
  loginBtn.style.color = "white";
  loginBtn.innerHTML =
    '<i class="fa fa-spinner fa-spin"></i> Logging you in...';
  loginBtn.disabled = true;

  const loginData = JSON.stringify({
    password,
    email,
  });
  const validationResult = validateLoginInput(password);

  if (Object.keys(validationResult).length > 0) {
    /* loginBtn.style.background = '#d64541'; */
    loginBtn.innerText = "Login";
    loginBtn.disabled = false;
    showErrorMessages(validationResult, errorDiv, errorMessage);
    return;
  }

  const url = "api/v1/auth/login";
  const resp = await fetch(url, {
    method: "POST",
    mode: "cors",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: loginData,
  });

  const result = await resp.json();
  const { statusCode, data, errors } = result;

  if (statusCode === 200) {
    localStorage.setItem("userToken", data.user.authToken);
    window.location.href = "chat.html";
  } else {
    showErrorMessages(errors, errorDiv, errorMessage);
    //loginBtn.style.background = '#d64541';
    loginBtn.innerText = "Login";
    loginBtn.disabled = false;
  }
};

loginForm.addEventListener("submit", login);

const signup = async (evt) => {
  evt.preventDefault();
  const email = document.getElementById("s-email").value;
  const password = document.getElementById("s-password").value;
  const confirmPassword = document.getElementById("re-password").value;
  const signupBtn = document.getElementById("signup-btn");
  /* loginBtn.style.background = 'darkgray'; */
  signupBtn.style.color = "white";
  signupBtn.innerHTML =
    '<i class="fa fa-spinner fa-spin"></i> Signing you up...';
  signupBtn.disabled = true;

  const signupData = JSON.stringify({
    password,
    email,
    confirmPassword,
  });
  const validationResult = validateSignupInput(password, confirmPassword);

  if (Object.keys(validationResult).length > 0) {
    /* loginBtn.style.background = '#d64541'; */
    signupBtn.innerText = "Sign Up";
    signupBtn.disabled = false;
    showErrorMessages(validationResult, signupErrorDiv, signupErrorMessage);
    return;
  }

  const url = "api/v1/auth/signup";
  const resp = await fetch(url, {
    method: "POST",
    mode: "cors",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: signupData,
  });

  const result = await resp.json();
  const { statusCode, data, errors } = result;

  if (statusCode === 201) {
    localStorage.setItem("userToken", data.user.authToken);
    window.location.href = "chat.html";
  } else {
    showErrorMessages(errors, signupErrorDiv, signupErrorMessage);
    //loginBtn.style.background = '#d64541';
    signupBtn.innerText = "Sign Up";
    signupBtn.disabled = false;
  }
};

signupForm.addEventListener("submit", signup);

window.onload = () => {
  try {
    const token = localStorage.getItem("userToken");
    if (token) {
      const { sub } = jwt_decode(token);

      if (sub) {
        window.location.href = "chat.html";
      }
    }
  } catch (error) {
    localStorage.removeItem("userToken");
  }
};
