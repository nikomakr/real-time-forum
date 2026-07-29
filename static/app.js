// RENTFORUM — Single Page Application

// PAGE REFERENCES
const loginPage = document.getElementById("login-page");
const registerPage = document.getElementById("register-page");
const forumPage = document.getElementById("forum-page");

// FORM REFERENCES
const loginForm = document.getElementById("loginForm");
const registerForm = document.getElementById("registerForm");

// LOGIN INPUTS
const loginUsername = document.getElementById("login-username");
const loginPassword = document.getElementById("login-password");

// REGISTER INPUTS
const nickname = document.getElementById("nickname");
const firstName = document.getElementById("first-name");
const lastName = document.getElementById("last-name");
const email = document.getElementById("email");
const password = document.getElementById("password");
const age = document.getElementById("age");
const gender = document.getElementById("gender");

// BUTTONS
const logoutButton = document.getElementById("logoutButton");
const showRegister = document.getElementById("show-register");
const showLogin = document.getElementById("show-login");

// ERROR CONTAINERS
let loginError;
let registerError;

// CREATE ERROR ELEMENTS
function createErrorContainers() {
  loginError = document.createElement("p");
  loginError.className = "form-error";
  loginForm.appendChild(loginError);
  registerError = document.createElement("p");
  registerError.className = "form-error";
  registerForm.appendChild(registerError);
}

// SHOW PAGE
function showPage(page) {
  document.querySelectorAll(".page").forEach((section) => {
    section.classList.remove("active");
    section.classList.add("hidden");
  });
  page.classList.remove("hidden");
  page.classList.add("active");
}

// SHOW LOGIN
function goToLogin() {
  clearErrors();
  loginForm.reset();
  showPage(loginPage);
}

// SHOW REGISTER
function goToRegister() {
  clearErrors();
  registerForm.reset();
  showPage(registerPage);
}

// SHOW FORUM
function goToForum() {
  clearErrors();
  showPage(forumPage);
}

// CLEAR ERRORS
function clearErrors() {
  if (loginError) {
    loginError.textContent = "";
  }
  if (registerError) {
    registerError.textContent = "";
  }
}

// CHECK SESSION
async function restoreSession() {
  try {
    const response = await fetch("/api/me", {
      credentials: "include",
    });
    if (response.ok) {
      goToForum();
    } else {
      goToLogin();
    }
  } catch (error) {
    console.error(error);
    goToLogin();
  }
}

// LINK NAVIGATION
showRegister.addEventListener("click", (event) => {
  event.preventDefault();
  goToRegister();
});

showLogin.addEventListener("click", (event) => {
  event.preventDefault();
  goToLogin();
});

// VALIDATION HELPERS
function showLoginError(message) {
  loginError.textContent = message;
}

function showRegisterError(message) {
  registerError.textContent = message;
}

// REGISTER VALIDATION
function validateRegistration() {
  clearErrors();
  if (nickname.value.trim() === "") {
    showRegisterError("Nickname is required.");
    return false;
  }
  if (nickname.value.includes("@")) {
    showRegisterError("Nickname cannot contain '@'.");
    return false;
  }
  if (firstName.value.trim() === "") {
    showRegisterError("First name is required.");
    return false;
  }
  if (lastName.value.trim() === "") {
    showRegisterError("Last name is required.");
    return false;
  }
  if (email.value.trim() === "") {
    showRegisterError("Email is required.");
    return false;
  }
  if (!email.value.includes("@")) {
    showRegisterError("Please enter a valid email.");
    return false;
  }
  if (password.value.trim() === "") {
    showRegisterError("Password is required.");
    return false;
  }
  const ageValue = Number(age.value);
  if (Number.isNaN(ageValue) || ageValue <= 0) {
    showRegisterError("Please enter a valid age.");
    return false;
  }
  if (gender.value === "") {
    showRegisterError("Please select a gender.");
    return false;
  }
  return true;
}

// REGISTER USER
async function registerUser(event) {
  event.preventDefault();
  if (!validateRegistration()) {
    return;
  }
  const payload = {
    nickname: nickname.value.trim(),
    first_name: firstName.value.trim(),
    last_name: lastName.value.trim(),
    email: email.value.trim(),
    password: password.value,
    age: Number(age.value),
    gender: gender.value,
  };
  try {
    const response = await fetch("/api/register", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });
    if (!response.ok) {
      let errorMessage = "Registration failed.";
      try {
        const error = await response.json();
        errorMessage = error.error || error.message || errorMessage;
      } catch (_) {}
      showRegisterError(errorMessage);
      return;
    }
    registerForm.reset();
    clearErrors();
    alert("Registration successful. Please log in.");
    goToLogin();
  } catch (error) {
    console.error(error);
    showRegisterError("Unable to connect to the server.");
  }
}

// REGISTER EVENT
registerForm.addEventListener("submit", registerUser);

// LOGIN VALIDATION
function validateLogin() {
  clearErrors();
  if (loginUsername.value.trim() === "") {
    showLoginError("Please enter your email or nickname.");
    return false;
  }
  if (loginPassword.value.trim() === "") {
    showLoginError("Please enter your password.");
    return false;
  }
  return true;
}

// LOGIN USER
async function loginUser(event) {
  event.preventDefault();
  if (!validateLogin()) {
    return;
  }
  const payload = {
    identifier: loginUsername.value.trim(),
    password: loginPassword.value,
  };
  try {
    const response = await fetch("/api/login", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });
    if (!response.ok) {
      let errorMessage = "Login failed.";
      try {
        const error = await response.json();
        errorMessage = error.error || error.message || errorMessage;
      } catch (_) {}
      showLoginError(errorMessage);
      return;
    }
    loginForm.reset();
    clearErrors();
    goToForum();
  } catch (error) {
    console.error(error);
    showLoginError("Unable to connect to the server.");
  }
}

// LOGIN EVENT
loginForm.addEventListener("submit", loginUser);

// LOGOUT
async function logoutUser() {
  try {
    const response = await fetch("/api/logout", {
      method: "POST",
      credentials: "include",
    });
    if (!response.ok) {
      console.error("Logout failed.");
    }
  } catch (error) {
    console.error(error);
  }
  loginForm.reset();
  registerForm.reset();
  clearErrors();
  goToLogin();
}

// LOGOUT EVENT
if (logoutButton) {
  logoutButton.addEventListener("click", logoutUser);
}

// AUTH GUARD
async function requireAuthentication() {
  try {
    const response = await fetch("/api/me", {
      credentials: "include",
    });
    if (!response.ok) {
      goToLogin();
      return false;
    }
    return true;
  } catch (error) {
    console.error(error);
    goToLogin();
    return false;
  }
}

// OPEN FORUM
async function openForum() {
  const authenticated = await requireAuthentication();
  if (!authenticated) {
    return;
  }
  goToForum();
}

// INITIALISE APP
async function initialiseApp() {
  createErrorContainers();
  await restoreSession();
}

// DOM READY
document.addEventListener("DOMContentLoaded", initialiseApp);

// DEBUG HELPERS
window.RentForum = {
  goToLogin,
  goToRegister,
  goToForum,
  openForum,
  restoreSession,
  logoutUser,
};
