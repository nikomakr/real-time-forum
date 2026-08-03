// RENTFORUM — Single Page Application

// =====================================================
// PAGE REFERENCES
// =====================================================
const loginPage = document.getElementById("login-page");
const registerPage = document.getElementById("register-page");
const forumPage = document.getElementById("forum-page");

// =====================================================
// FORM REFERENCES
// =====================================================
const loginForm = document.getElementById("loginForm");
const registerForm = document.getElementById("registerForm");
const discussionForm = document.getElementById("discussionForm");
const commentForm = document.getElementById("commentForm");

// =====================================================
// LOGIN INPUTS
// =====================================================
const loginUsername = document.getElementById("login-username");
const loginPassword = document.getElementById("login-password");

// =====================================================
// REGISTER INPUTS
// =====================================================
const nickname = document.getElementById("nickname");
const firstName = document.getElementById("first-name");
const lastName = document.getElementById("last-name");
const email = document.getElementById("email");
const password = document.getElementById("password");
const age = document.getElementById("age");
const gender = document.getElementById("gender");

// =====================================================
// BUTTONS
// =====================================================
const logoutButton = document.getElementById("logoutButton");
const showRegister = document.getElementById("show-register");
const showLogin = document.getElementById("show-login");
const newPostButton = document.getElementById("new-post");

// =====================================================
// FORUM ELEMENTS
// =====================================================
const discussionList = document.getElementById("discussionList");
const createDiscussionModal = document.getElementById("createDiscussionModal");
const discussionModal = document.getElementById("discussionModal");
const commentsContainer = document.getElementById("commentsContainer");
const discussionTemplate = document.getElementById("discussionTemplate");
const commentTemplate = document.getElementById("commentTemplate");

// =====================================================
// STATE
// =====================================================
let currentPostId = null;
let activeCategory = "all";

// =====================================================
// ERROR CONTAINERS
// =====================================================
let loginError;
let registerError;

// =====================================================
// CREATE ERROR ELEMENTS
// =====================================================
function createErrorContainers() {
  loginError = document.createElement("p");
  loginError.className = "form-error";
  loginForm.appendChild(loginError);

  registerError = document.createElement("p");
  registerError.className = "form-error";
  registerForm.appendChild(registerError);
}

// =====================================================
// SHOW PAGE
// =====================================================
function showPage(page) {
  document.querySelectorAll(".page").forEach((section) => {
    section.classList.remove("active");
    section.classList.add("hidden");
  });
  page.classList.remove("hidden");
  page.classList.add("active");
}

// =====================================================
// NAVIGATION
// =====================================================
function goToLogin() {
  clearErrors();
  loginForm.reset();
  showPage(loginPage);
}

function goToRegister() {
  clearErrors();
  registerForm.reset();
  showPage(registerPage);
}

async function goToForum() {
  clearErrors();
  showPage(forumPage);
  await Promise.all([loadFeed(activeCategory), loadCategoryTiles()]);
}

// =====================================================
// CLEAR ERRORS
// =====================================================
function clearErrors() {
  if (loginError) loginError.textContent = "";
  if (registerError) registerError.textContent = "";
}

// =====================================================
// SESSION RESTORE
// =====================================================
async function restoreSession() {
  try {
    const response = await fetch("/api/me", { credentials: "include" });
    if (response.ok) {
      await goToForum();
    } else {
      goToLogin();
    }
  } catch (error) {
    console.error(error);
    goToLogin();
  }
}

// =====================================================
// LINK NAVIGATION
// =====================================================
showRegister.addEventListener("click", (event) => {
  event.preventDefault();
  goToRegister();
});

showLogin.addEventListener("click", (event) => {
  event.preventDefault();
  goToLogin();
});

// =====================================================
// VALIDATION HELPERS
// =====================================================
function showLoginError(message) {
  loginError.textContent = message;
}

function showRegisterError(message) {
  registerError.textContent = message;
}

// =====================================================
// REGISTER VALIDATION
// =====================================================
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

// =====================================================
// REGISTER USER
// =====================================================
async function registerUser(event) {
  event.preventDefault();
  if (!validateRegistration()) return;

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
      headers: { "Content-Type": "application/json" },
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

registerForm.addEventListener("submit", registerUser);

// =====================================================
// LOGIN VALIDATION
// =====================================================
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

// =====================================================
// LOGIN USER
// =====================================================
async function loginUser(event) {
  event.preventDefault();
  if (!validateLogin()) return;

  const payload = {
    identifier: loginUsername.value.trim(),
    password: loginPassword.value,
  };

  try {
    const response = await fetch("/api/login", {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
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
    await goToForum();
  } catch (error) {
    console.error(error);
    showLoginError("Unable to connect to the server.");
  }
}

loginForm.addEventListener("submit", loginUser);

// =====================================================
// LOGOUT
// =====================================================
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
  activeCategory = "all";
  goToLogin();
}

if (logoutButton) {
  logoutButton.addEventListener("click", logoutUser);
}

// =====================================================
// AUTH GUARD
// =====================================================
async function requireAuthentication() {
  try {
    const response = await fetch("/api/me", { credentials: "include" });
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

// =====================================================
// UTILITY HELPERS
// =====================================================
function escapeHTML(str) {
  return String(str)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

function truncate(str, max) {
  if (!str) return "";
  return str.length > max ? str.slice(0, max) + "…" : str;
}

function formatDate(dateStr) {
  if (!dateStr) return "";
  const date = new Date(dateStr);
  if (isNaN(date)) return dateStr;
  return date.toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function slugify(str) {
  if (!str) return "";
  return str
    .toLowerCase()
    .replace(/\s+/g, "-")
    .replace(/[^a-z0-9-]/g, "");
}

function splitAndTrim(str, sep) {
  if (!str) return [];
  return str
    .split(sep)
    .map((s) => s.trim())
    .filter((s) => s !== "");
}

// =====================================================
// CATEGORY FILTER — group tiles + category picker modal
// Tiles are rendered from GET /api/categories (category_groups
// table); clicking a group tile opens a picker listing that
// group's categories (categories table) to filter the feed by.
// =====================================================
const categoryTiles = document.getElementById("categoryTiles");
const categoryPickerModal = document.getElementById("categoryPickerModal");
const categoryPickerList = document.getElementById("categoryPickerList");
const categoryPickerHeading = document.getElementById("categoryPickerHeading");
const allCategoriesTile = document.getElementById("category-all");

let categoryGroups = [];

function setActiveCategoryTile(tile) {
  document.querySelectorAll(".category-tile").forEach((el) => {
    el.classList.remove("active-category");
  });
  if (tile) tile.classList.add("active-category");
}

function selectCategory(name, tile) {
  activeCategory = name;
  setActiveCategoryTile(tile);
  closeCategoryPicker();
  loadFeed(activeCategory);
}

function openCategoryPicker(group, tile) {
  if (!categoryPickerModal || !categoryPickerList) return;
  categoryPickerHeading.textContent = group.name;
  categoryPickerList.innerHTML = "";

  if (!group.categories || group.categories.length === 0) {
    const empty = document.createElement("li");
    empty.textContent = "No categories in this group yet.";
    categoryPickerList.appendChild(empty);
  } else {
    let lastSubGroup;
    group.categories.forEach((cat) => {
      if (cat.sub_group && cat.sub_group !== lastSubGroup) {
        const label = document.createElement("li");
        label.className = "sub-group-label";
        label.textContent = cat.sub_group;
        categoryPickerList.appendChild(label);
        lastSubGroup = cat.sub_group;
      }
      const item = document.createElement("li");
      item.textContent = cat.name;
      item.setAttribute("role", "listitem");
      item.tabIndex = 0;
      if (activeCategory === cat.name) item.classList.add("active-category");
      item.addEventListener("click", () => selectCategory(cat.name, tile));
      categoryPickerList.appendChild(item);
    });
  }

  categoryPickerModal.classList.remove("hidden");
}

function closeCategoryPicker() {
  if (categoryPickerModal) categoryPickerModal.classList.add("hidden");
}

function renderCategoryTiles() {
  if (!categoryTiles) return;

  categoryTiles.querySelectorAll(".category-tile[data-group-id]").forEach((el) => el.remove());

  categoryGroups.forEach((group) => {
    const tile = document.createElement("div");
    tile.className = "category-tile";
    tile.dataset.groupId = group.id;
    tile.setAttribute("role", "listitem");
    tile.tabIndex = 0;
    tile.textContent = group.name;
    tile.addEventListener("click", () => openCategoryPicker(group, tile));
    categoryTiles.appendChild(tile);
  });
}

async function loadCategoryTiles() {
  try {
    const response = await fetch("/api/categories", { credentials: "include" });
    if (!response.ok) return;
    categoryGroups = await response.json();
    renderCategoryTiles();
  } catch (error) {
    console.error(error);
  }
}

function initCategoryFilter() {
  if (allCategoriesTile) {
    allCategoriesTile.addEventListener("click", () => selectCategory("all", allCategoriesTile));
  }
}

// =====================================================
// FEED VIEW
// =====================================================
async function loadFeed(category = "all") {
  if (!discussionList) return;
  discussionList.innerHTML = '<p class="loading-text">Loading discussions…</p>';

  let url = "/api/posts";
  if (category && category !== "all") {
    url += `?category=${encodeURIComponent(category)}`;
  }

  try {
    const response = await fetch(url, { credentials: "include" });
    if (!response.ok) {
      discussionList.innerHTML =
        '<p class="feed-error">Could not load discussions.</p>';
      return;
    }
    const posts = await response.json();
    renderFeed(posts);
  } catch (error) {
    console.error(error);
    discussionList.innerHTML =
      '<p class="feed-error">Could not connect to the server.</p>';
  }
}

function renderFeed(posts) {
  if (!discussionList) return;
  discussionList.innerHTML = "";

  if (!posts || posts.length === 0) {
    discussionList.innerHTML =
      '<p class="no-results">No discussions found. Be the first to post!</p>';
    return;
  }

  posts.forEach((post) => {
    const card = buildPostCard(post);
    discussionList.appendChild(card);
  });
}

// =====================================================
// POST CARD COMPONENT
// FIX: backend returns post.categories (array), not post.category (string)
// =====================================================
function buildPostCard(post) {
  const template = discussionTemplate.content.cloneNode(true);

  // Use first category for badge display — backend returns array
  const firstCategory =
    Array.isArray(post.categories) && post.categories.length > 0
      ? post.categories[0]
      : "";

  const categoryEl = template.querySelector(".category");
  categoryEl.textContent = firstCategory;
  categoryEl.className = `category ${slugify(firstCategory)}`;

  template.querySelector("h3").textContent = post.title || "";

  template.querySelector(".discussion-meta").innerHTML =
    `<span>Posted by ${escapeHTML(post.author || "Anonymous")}</span>
     <span>•</span>
     <span>${formatDate(post.created_at)}</span>`;

  template.querySelector(".discussion-preview").textContent = truncate(
    post.content || "",
    160,
  );

  const readBtn = template.querySelector(".read-btn");
  readBtn.dataset.postId = post.id;
  readBtn.addEventListener("click", () => openPostDetail(post.id));

  const commentBtn = template.querySelector(".comment-btn");
  commentBtn.textContent = `💬 ${post.comment_count || 0}`;
  commentBtn.dataset.postId = post.id;
  commentBtn.addEventListener("click", () => openPostDetail(post.id));

  const likeBtn = template.querySelector(".like-btn");
  if (likeBtn) likeBtn.textContent = `👍 ${post.like_count || 0}`;

  return template;
}

// =====================================================
// POST DETAIL VIEW
// FIX: backend returns post.categories (array), not post.category (string)
// =====================================================
async function openPostDetail(postId) {
  currentPostId = postId;

  try {
    const [postRes, commentsRes] = await Promise.all([
      fetch(`/api/posts/${postId}`, { credentials: "include" }),
      fetch(`/api/posts/${postId}/comments`, { credentials: "include" }),
    ]);

    if (!postRes.ok) {
      alert("Could not load this discussion.");
      return;
    }

    const post = await postRes.json();
    const comments = commentsRes.ok ? await commentsRes.json() : [];

    renderPostDetail(post);
    renderComments(comments);
    discussionModal.classList.remove("hidden");
  } catch (error) {
    console.error(error);
    alert("Could not connect to the server.");
  }
}

function renderPostDetail(post) {
  // Use first category for badge — backend returns array
  const firstCategory =
    Array.isArray(post.categories) && post.categories.length > 0
      ? post.categories[0]
      : "";

  const categoryEl = discussionModal.querySelector(".category");
  categoryEl.textContent = firstCategory;
  categoryEl.className = `category ${slugify(firstCategory)}`;

  document.getElementById("discussionHeading").textContent = post.title || "";

  discussionModal.querySelector(".discussion-meta").innerHTML =
    `<span>Posted by ${escapeHTML(post.author || "Anonymous")}</span>
     <span>•</span>
     <span>${formatDate(post.created_at)}</span>`;

  discussionModal.querySelector(".discussion-content p").textContent =
    post.content || "";
}

// =====================================================
// COMMENTS LIST
// =====================================================
function renderComments(comments) {
  if (!commentsContainer) return;
  commentsContainer.innerHTML = "";

  if (!comments || comments.length === 0) {
    commentsContainer.innerHTML =
      '<p class="no-results">No replies yet. Be the first!</p>';
    return;
  }

  comments.forEach((comment) => {
    const template = commentTemplate.content.cloneNode(true);
    template.querySelector("strong").textContent =
      comment.author || "Anonymous";
    template.querySelector("p").textContent = comment.content || "";
    template.querySelector("small").textContent = formatDate(
      comment.created_at,
    );
    commentsContainer.appendChild(template);
  });
}

// =====================================================
// COMMENT FORM
// =====================================================
async function submitComment(event) {
  event.preventDefault();
  if (!currentPostId) return;

  const commentText = document.getElementById("commentText");
  const content = commentText.value.trim();
  if (!content) return;

  try {
    const response = await fetch(`/api/posts/${currentPostId}/comments`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ content }),
    });

    if (!response.ok) {
      alert("Could not post reply.");
      return;
    }

    commentText.value = "";

    // Reload comments after posting
    const commentsRes = await fetch(`/api/posts/${currentPostId}/comments`, {
      credentials: "include",
    });
    const comments = commentsRes.ok ? await commentsRes.json() : [];
    renderComments(comments);
  } catch (error) {
    console.error(error);
    alert("Could not connect to the server.");
  }
}

if (commentForm) {
  commentForm.addEventListener("submit", submitComment);
}

// =====================================================
// CREATE POST VIEW
// FIX: backend expects categories as array of IDs, not a single category name
// The select value is the category ID — see index.html option values
// =====================================================
function openCreatePost() {
  if (createDiscussionModal) {
    createDiscussionModal.classList.remove("hidden");
  }
}

function closeCreatePost() {
  if (createDiscussionModal) {
    createDiscussionModal.classList.add("hidden");
  }
  if (discussionForm) discussionForm.reset();
}

async function submitPost(event) {
  event.preventDefault();

  const titleEl = document.getElementById("discussionTitle");
  const categoryEl = document.getElementById("discussionCategory");
  const bodyEl = document.getElementById("discussionBody");

  const title = titleEl ? titleEl.value.trim() : "";
  // category value is now the category ID from the select option value
  const categoryId = categoryEl ? categoryEl.value : "";
  const content = bodyEl ? bodyEl.value.trim() : "";

  if (!title) {
    alert("Please enter a title.");
    return;
  }
  if (!categoryId) {
    alert("Please select a category.");
    return;
  }
  if (!content) {
    alert("Please enter the discussion content.");
    return;
  }

  try {
    const response = await fetch("/api/posts", {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      // categories must be an array of category IDs
      body: JSON.stringify({ title, content, categories: [categoryId] }),
    });

    if (!response.ok) {
      let errorMessage = "Could not create discussion.";
      try {
        const error = await response.json();
        errorMessage = error.error || error.message || errorMessage;
      } catch (_) {}
      alert(errorMessage);
      return;
    }

    closeCreatePost();
    await loadFeed(activeCategory);
  } catch (error) {
    console.error(error);
    alert("Could not connect to the server.");
  }
}

if (discussionForm) {
  discussionForm.addEventListener("submit", submitPost);
}

if (newPostButton) {
  newPostButton.addEventListener("click", openCreatePost);
}

// =====================================================
// CLOSE MODALS
// =====================================================
document.querySelectorAll(".close-modal").forEach((btn) => {
  btn.addEventListener("click", () => {
    document.querySelectorAll(".modal").forEach((modal) => {
      modal.classList.add("hidden");
    });
    currentPostId = null;
  });
});

// =====================================================
// INITIALISE APP
// =====================================================
async function initialiseApp() {
  createErrorContainers();
  initCategoryFilter();
  await restoreSession();
}

document.addEventListener("DOMContentLoaded", initialiseApp);

// =====================================================
// DEBUG HELPERS
// =====================================================
window.RentForum = {
  goToLogin,
  goToRegister,
  goToForum,
  restoreSession,
  logoutUser,
  loadFeed,
  openPostDetail,
};
