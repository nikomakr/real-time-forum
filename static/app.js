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
const logoHomeButton = document.getElementById("logoHomeButton");
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
let categoryGroups = [];
let myUserID = null;

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
  disconnectWebSocket();
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
  await fetchMyUserID();
  connectWebSocket();
  await Promise.all([
    loadFeed(activeCategory),
    loadCategoryTiles(),
    loadUsers(),
    loadProfile(),
  ]);
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
      const data = await response.json();
      myUserID = data.user_id;
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
// FETCH MY USER ID
// =====================================================
async function fetchMyUserID() {
  try {
    const response = await fetch("/api/me", { credentials: "include" });
    if (response.ok) {
      const data = await response.json();
      myUserID = data.user_id;
    }
  } catch (error) {
    console.error("[fetchMyUserID]", error);
  }
}

// =====================================================
// PROFILE
// =====================================================
async function loadProfile() {
  try {
    const response = await fetch("/api/profile", { credentials: "include" });
    if (!response.ok) return;
    const profile = await response.json();

    const usernameEl = document.getElementById("username");
    if (usernameEl) usernameEl.textContent = profile.nickname;

    const nameEl = document.getElementById("profileName");
    if (nameEl) nameEl.textContent = profile.nickname;

    const emailEl = document.getElementById("profileEmail");
    if (emailEl) emailEl.textContent = profile.email;

    const postsEl = document.getElementById("profileStatPosts");
    if (postsEl) postsEl.textContent = profile.post_count;

    const commentsEl = document.getElementById("profileStatComments");
    if (commentsEl) commentsEl.textContent = profile.comment_count;

    const likesEl = document.getElementById("profileStatLikes");
    if (likesEl) likesEl.textContent = profile.like_count;
  } catch (error) {
    console.error("[loadProfile]", error);
  }
}

function openProfileModal() {
  const dropdown = document.getElementById("profileDropdown");
  if (dropdown) dropdown.classList.add("hidden");

  const modal = document.getElementById("profileModal");
  if (modal) modal.classList.remove("hidden");

  loadProfile();
}

const profileMenuButton = document.getElementById("profileMenuButton");
if (profileMenuButton) {
  profileMenuButton.addEventListener("click", () => {
    const dropdown = document.getElementById("profileDropdown");
    if (dropdown) {
      const isHidden = dropdown.classList.contains("hidden");
      dropdown.classList.toggle("hidden");
      profileMenuButton.setAttribute(
        "aria-expanded",
        isHidden ? "true" : "false",
      );
    }
  });
}

const openProfileButton = document.getElementById("openProfileButton");
if (openProfileButton) {
  openProfileButton.addEventListener("click", openProfileModal);
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
  if (password.value.length < 12) {
    showRegisterError("Password must be at least 12 characters long.");
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
    await fetch("/api/logout", { method: "POST", credentials: "include" });
  } catch (error) {
    console.error(error);
  }
  disconnectWebSocket();
  myUserID = null;
  currentChatUserID = null;
  currentChatUserName = null;
  loginForm.reset();
  registerForm.reset();
  clearErrors();
  activeCategory = "all";
  goToLogin();
}

if (logoutButton) {
  logoutButton.addEventListener("click", logoutUser);
}

if (logoHomeButton) {
  logoHomeButton.addEventListener("click", () => {
    window.scrollTo({ top: 0, behavior: "smooth" });
  });
}

const logoutDropdown = document.getElementById("logoutDropdown");
if (logoutDropdown) {
  logoutDropdown.addEventListener("click", logoutUser);
}

const profileLogoutButton = document.getElementById("profileLogoutButton");
if (profileLogoutButton) {
  profileLogoutButton.addEventListener("click", logoutUser);
}

const logoutFromDiscussion = document.getElementById("logoutFromDiscussion");
if (logoutFromDiscussion) {
  logoutFromDiscussion.addEventListener("click", logoutUser);
}

const logoutFromMessages = document.getElementById("logoutFromMessages");
if (logoutFromMessages) {
  logoutFromMessages.addEventListener("click", logoutUser);
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
// CATEGORY TILES — loaded from GET /api/categories
// =====================================================
async function loadCategoryTiles() {
  try {
    const response = await fetch("/api/categories", { credentials: "include" });
    if (!response.ok) return;
    categoryGroups = await response.json();
    renderCategoryTiles();
  } catch (error) {
    console.error("[loadCategoryTiles]", error);
  }
}

function renderCategoryTiles() {
  const tilesContainer = document.getElementById("categoryTiles");
  if (!tilesContainer) return;

  // Keep the "All Discussions" tile, remove any previously injected group tiles
  const existingGroupTiles = tilesContainer.querySelectorAll(
    ".category-tile[data-group]",
  );
  existingGroupTiles.forEach((t) => t.remove());

  categoryGroups.forEach((group) => {
    const tile = document.createElement("div");
    tile.className = "category-tile";
    tile.dataset.group = group.id;
    tile.textContent = group.name;
    tile.setAttribute("role", "button");
    tile.setAttribute("tabindex", "0");
    tile.setAttribute("aria-label", `Browse ${group.name} categories`);
    tile.addEventListener("click", () => openCategoryPicker(group, tile));
    tile.addEventListener("keydown", (e) => {
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        openCategoryPicker(group, tile);
      }
    });
    tilesContainer.appendChild(tile);
  });
}

function openCategoryPicker(group, tile) {
  const modal = document.getElementById("categoryPickerModal");
  const heading = document.getElementById("categoryPickerHeading");
  const list = document.getElementById("categoryPickerList");
  if (!modal || !list || !heading) return;

  heading.textContent = group.name;
  list.innerHTML = "";

  let currentSubGroup = null;

  group.categories.forEach((cat) => {
    if (cat.sub_group && cat.sub_group !== currentSubGroup) {
      currentSubGroup = cat.sub_group;
      const label = document.createElement("li");
      label.className = "sub-group-label";
      label.textContent = cat.sub_group;
      list.appendChild(label);
    }

    const li = document.createElement("li");
    li.textContent = cat.name;
    li.setAttribute("role", "button");
    li.setAttribute("tabindex", "0");
    if (cat.name === activeCategory) li.classList.add("active-category");
    li.addEventListener("click", () => selectCategory(cat.name, tile));
    li.addEventListener("keydown", (e) => {
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        selectCategory(cat.name, tile);
      }
    });
    list.appendChild(li);
  });

  modal.classList.remove("hidden");
}

function closeCategoryPicker() {
  const modal = document.getElementById("categoryPickerModal");
  if (modal) modal.classList.add("hidden");
}

function setActiveCategoryTile(tile) {
  document
    .querySelectorAll(".category-tile")
    .forEach((t) => t.classList.remove("active-category"));
  if (tile) tile.classList.add("active-category");
}

function selectCategory(name, tile) {
  activeCategory = name;
  setActiveCategoryTile(tile);
  closeCategoryPicker();
  loadFeed(activeCategory);
}

function initCategoryFilter() {
  const allTile = document.getElementById("category-all");
  if (allTile) {
    allTile.addEventListener("click", () => {
      activeCategory = "all";
      setActiveCategoryTile(allTile);
      loadFeed("all");
    });
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
// =====================================================
function buildPostCard(post) {
  const template = discussionTemplate.content.cloneNode(true);

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
  commentBtn.dataset.postId = post.id;
  const commentCount = commentBtn.querySelector
    ? commentBtn
    : template.querySelector(".comment-btn");
  const countSpan = document.createTextNode(` ${post.comment_count || 0}`);
  commentBtn.appendChild(countSpan);
  commentBtn.addEventListener("click", () => openPostDetail(post.id));

  const likeBtn = template.querySelector(".like-btn");
  if (likeBtn) {
    likeBtn.dataset.postId = post.id;
    setLikeButtonState(likeBtn, post.liked, post.like_count || 0);
    likeBtn.addEventListener("click", () => toggleLike(post.id, likeBtn));
  }

  return template;
}

// =====================================================
// LIKE BUTTON
// =====================================================
function setLikeButtonState(likeBtn, liked, count) {
  likeBtn.classList.toggle("liked", !!liked);
  likeBtn.setAttribute("aria-pressed", liked ? "true" : "false");
  const countEl = likeBtn.querySelector(".like-count");
  if (countEl) countEl.textContent = ` ${count}`;
}

async function toggleLike(postId, likeBtn) {
  likeBtn.disabled = true;
  try {
    const response = await fetch(`/api/posts/${postId}/like`, {
      method: "POST",
      credentials: "include",
    });
    if (!response.ok) return;
    const { liked, like_count } = await response.json();
    setLikeButtonState(likeBtn, liked, like_count);
  } catch (error) {
    console.error(error);
  } finally {
    likeBtn.disabled = false;
  }
}

// =====================================================
// POST DETAIL VIEW
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
// =====================================================
function openCreatePost() {
  if (createDiscussionModal) createDiscussionModal.classList.remove("hidden");
}

function closeCreatePost() {
  if (createDiscussionModal) createDiscussionModal.classList.add("hidden");
  if (discussionForm) discussionForm.reset();
}

async function submitPost(event) {
  event.preventDefault();

  const titleEl = document.getElementById("discussionTitle");
  const categoryEl = document.getElementById("discussionCategory");
  const bodyEl = document.getElementById("discussionBody");

  const title = titleEl ? titleEl.value.trim() : "";
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
    const picker = document.getElementById("categoryPickerModal");
    if (picker) picker.classList.add("hidden");
    currentPostId = null;
    resetChatPanel();
  });
});

// Close a modal when clicking its backdrop (outside .modal-content)
document.querySelectorAll(".modal").forEach((modal) => {
  modal.addEventListener("click", (event) => {
    if (event.target !== modal) return;
    modal.classList.add("hidden");
    currentPostId = null;
    if (modal.id === "messagesModal") resetChatPanel();
  });
});

// =====================================================
// WEBSOCKET CLIENT — Epic 8
// =====================================================
let socket = null;
let currentChatUserID = null;
let currentChatUserName = null;
let messageOffset = 0;
let isLoadingMessages = false;
let hasMoreMessages = true;
let selectedImageFile = null;

function connectWebSocket() {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const wsURL = `${protocol}//${window.location.host}/ws`;

  socket = new WebSocket(wsURL);

  socket.addEventListener("open", () => {
    // #socketStatus shows the presence of the person you're chatting with,
    // not our own connection state — see updateChatStatus().
  });

  socket.addEventListener("message", (event) => {
    // Server may batch messages separated by newlines
    const lines = event.data.split("\n").filter((l) => l.trim());
    lines.forEach((line) => {
      try {
        const envelope = JSON.parse(line);
        handleWebSocketMessage(envelope);
      } catch (e) {
        console.error("[WS parse error]", e);
      }
    });
  });

  socket.addEventListener("close", () => {
    // Reconnect after 3 seconds if still logged in
    if (myUserID) {
      setTimeout(connectWebSocket, 3000);
    }
  });

  socket.addEventListener("error", (e) => {
    console.error("[WS error]", e);
  });
}

function disconnectWebSocket() {
  if (socket) {
    socket.close();
    socket = null;
  }
}

function handleWebSocketMessage(envelope) {
  switch (envelope.type) {
    case "chat_message":
      handleIncomingChatMessage(envelope.payload);
      break;
    case "presence":
      handlePresenceUpdate(envelope.payload);
      break;
    case "notification":
      handleNotification(envelope.payload);
      break;
    default:
      console.warn("[WS unknown type]", envelope.type);
  }
}

function handleIncomingChatMessage(payload) {
  const isCurrentConversation =
    payload.sender_id === currentChatUserID ||
    payload.receiver_id === currentChatUserID;

  if (isCurrentConversation) {
    appendMessage(payload);
    scrollChatToBottom();
  }

  // Refresh user list so ordering updates
  loadUsers();
}

function handlePresenceUpdate(payload) {
  const onlineIDs = payload.online_user_ids || [];
  updateOnlineCounter(onlineIDs.length);
  loadUsers();
}

function handleNotification(payload) {
  // Only show notification if the message is not from the currently open conversation
  if (payload.sender_id !== currentChatUserID) {
    showNotificationBadge(payload);
    showToast(`New message from ${escapeHTML(payload.sender_name)}`);
  }
}

// =====================================================
// USER LIST — Epic 8
// =====================================================
async function loadUsers() {
  try {
    const response = await fetch("/api/users", { credentials: "include" });
    if (!response.ok) return;
    const users = await response.json();
    renderUserLists(users);
  } catch (error) {
    console.error("[loadUsers]", error);
  }
}

function renderUserLists(users) {
  const onlineList = document.getElementById("onlineList");
  const offlineList = document.getElementById("offlineList");
  if (!onlineList || !offlineList) return;

  onlineList.innerHTML = "";
  offlineList.innerHTML = "";

  const onlineUsers = users.filter((u) => u.is_online);
  const offlineUsers = users.filter((u) => !u.is_online);

  // Update panel headings with live counts
  const onlineHeading = onlineList.previousElementSibling;
  const offlineHeading = offlineList.previousElementSibling;
  if (onlineHeading)
    onlineHeading.textContent = `Online Members (${onlineUsers.length})`;
  if (offlineHeading)
    offlineHeading.textContent = `Offline (${offlineUsers.length})`;

  if (onlineUsers.length === 0) {
    const li = document.createElement("li");
    li.className = "user-list-empty";
    li.textContent = "No members online";
    onlineList.appendChild(li);
  } else {
    onlineUsers.forEach((user) =>
      onlineList.appendChild(buildUserListItem(user)),
    );
  }

  if (offlineUsers.length === 0) {
    const li = document.createElement("li");
    li.className = "user-list-empty";
    li.textContent = "No offline members";
    offlineList.appendChild(li);
  } else {
    offlineUsers.forEach((user) =>
      offlineList.appendChild(buildUserListItem(user)),
    );
  }

  updateOnlineCounter(onlineUsers.length);

  // Keep the open conversation's status pill in sync with live presence
  if (currentChatUserID) {
    const chatUser = users.find((u) => u.id === currentChatUserID);
    if (chatUser) updateChatStatus(chatUser.is_online);
  }
}

// Reflect whether the person you're chatting with (not you) is online,
// in the messages modal's status pill.
function updateChatStatus(isOnline) {
  const status = document.getElementById("socketStatus");
  if (!status) return;
  status.classList.toggle("online", isOnline);
  status.classList.toggle("offline", !isOnline);
  status.innerHTML = `<svg class="icon" aria-hidden="true" focusable="false"><use href="#icon-dot"/></svg> ${isOnline ? "Online" : "Offline"}`;
}

function buildUserListItem(user) {
  const li = document.createElement("li");
  li.className = `user-item${user.is_online ? " is-online" : " is-offline"}`;
  if (user.id === currentChatUserID) {
    li.classList.add("active-chat");
  }
  li.setAttribute("role", "button");
  li.setAttribute("tabindex", "0");
  li.setAttribute("aria-label", `Chat with ${escapeHTML(user.nickname)}`);

  li.innerHTML = `
    <span class="user-status-dot${user.is_online ? " online" : ""}"></span>
    <span class="user-nickname">${escapeHTML(user.nickname)}</span>
  `;

  li.addEventListener("click", () => openConversation(user));
  li.addEventListener("keydown", (e) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      openConversation(user);
    }
  });

  return li;
}

function updateOnlineCounter(count) {
  const counter = document.getElementById("onlineCounter");
  if (counter) counter.innerHTML = `Online: <strong>${count}</strong>`;
}

// =====================================================
// CONVERSATION — Epic 8
// =====================================================
async function openConversation(user) {
  currentChatUserID = user.id;
  currentChatUserName = user.nickname;
  messageOffset = 0;
  hasMoreMessages = true;

  // Open the private messages modal — plenty of room to read/write,
  // unlike the old cramped sidebar chat-window.
  const messagesModal = document.getElementById("messagesModal");
  if (messagesModal) messagesModal.classList.remove("hidden");

  // Update chat heading to show who you are chatting with
  const chatHeading = document.getElementById("chatHeading");
  if (chatHeading) chatHeading.textContent = user.nickname;

  // Show whether this member is currently online
  updateChatStatus(user.is_online);

  // Enable message form
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("sendMessageButton");
  const imageInput = document.getElementById("messageImageInput");
  const attachButton = document.getElementById("attachImageButton");
  if (messageInput) messageInput.disabled = false;
  if (sendButton) sendButton.disabled = false;
  if (imageInput) imageInput.disabled = false;
  if (attachButton) attachButton.disabled = false;

  // Clear chat window
  const chatWindow = document.getElementById("chatWindow");
  if (chatWindow) {
    chatWindow.innerHTML = '<p class="loading-text">Loading messages…</p>';
    chatWindow.removeEventListener("scroll", handleChatScroll);
    chatWindow.addEventListener("scroll", handleChatScroll);
  }

  // Highlight active user in list
  document
    .querySelectorAll(".user-item")
    .forEach((li) => li.classList.remove("active-chat"));

  await loadMessages(false);
  loadUsers();
}

// Reset the private chat back to its empty state when the messages modal
// is closed, so a re-opened conversation always starts fresh and incoming
// messages while it's closed are treated as "not the open conversation".
function resetChatPanel() {
  currentChatUserID = null;
  currentChatUserName = null;

  const chatHeading = document.getElementById("chatHeading");
  if (chatHeading) chatHeading.textContent = "Private Messages";

  const status = document.getElementById("socketStatus");
  if (status) {
    status.classList.remove("online", "offline");
    status.innerHTML = "";
  }

  const chatWindow = document.getElementById("chatWindow");
  if (chatWindow) {
    chatWindow.innerHTML =
      '<p class="no-results">Select a member to start chatting.</p>';
    chatWindow.removeEventListener("scroll", handleChatScroll);
  }

  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("sendMessageButton");
  const imageInput = document.getElementById("messageImageInput");
  const attachButton = document.getElementById("attachImageButton");
  if (messageInput) messageInput.disabled = true;
  if (sendButton) sendButton.disabled = true;
  if (imageInput) imageInput.disabled = true;
  if (attachButton) attachButton.disabled = true;

  document
    .querySelectorAll(".user-item")
    .forEach((li) => li.classList.remove("active-chat"));
}

async function loadMessages(prepend = false) {
  if (!currentChatUserID || isLoadingMessages) return;
  isLoadingMessages = true;

  const chatWindow = document.getElementById("chatWindow");

  try {
    const response = await fetch(
      `/api/messages/${currentChatUserID}?offset=${messageOffset}`,
      { credentials: "include" },
    );

    if (!response.ok) {
      if (chatWindow && !prepend) {
        chatWindow.innerHTML =
          '<p class="feed-error">Could not load messages.</p>';
      }
      return;
    }

    const messages = await response.json();
    if (!chatWindow) return;

    if (!prepend) {
      chatWindow.innerHTML = "";
    }

    if (!messages || messages.length === 0) {
      if (!prepend) {
        chatWindow.innerHTML =
          '<p class="no-results">No messages yet. Say hello!</p>';
      }
      hasMoreMessages = false;
      return;
    }

    if (messages.length < 10) {
      hasMoreMessages = false;
    }

    if (prepend) {
      // Save scroll height before prepending so position is preserved
      const prevScrollHeight = chatWindow.scrollHeight;
      // Insert as one fragment — inserting bubbles one-by-one before
      // firstChild would reverse their order within the batch.
      const fragment = document.createDocumentFragment();
      messages.forEach((msg) => {
        fragment.appendChild(buildMessageBubble(msg));
      });
      chatWindow.insertBefore(fragment, chatWindow.firstChild);
      chatWindow.scrollTop = chatWindow.scrollHeight - prevScrollHeight;
    } else {
      messages.forEach((msg) => {
        chatWindow.appendChild(buildMessageBubble(msg));
      });
      scrollChatToBottom();
    }

    messageOffset += messages.length;
  } catch (error) {
    console.error("[loadMessages]", error);
  } finally {
    isLoadingMessages = false;
  }
}

function messageImageHTML(imageUrl) {
  if (!imageUrl) return "";
  return `<img src="${escapeHTML(imageUrl)}" alt="Shared image" loading="lazy" />`;
}

function buildMessageBubble(msg) {
  const isSent = msg.sender_id === myUserID;
  const div = document.createElement("div");
  div.className = `message ${isSent ? "sent" : "received"}`;
  div.innerHTML = `
    <span class="sender">${escapeHTML(msg.sender_name)}</span>
    ${msg.content ? `<p>${escapeHTML(msg.content)}</p>` : ""}
    ${messageImageHTML(msg.image_url)}
    <small>${formatDate(msg.created_at)}</small>
  `;
  return div;
}

function appendMessage(payload) {
  const chatWindow = document.getElementById("chatWindow");
  if (!chatWindow) return;

  const noResults = chatWindow.querySelector(".no-results");
  if (noResults) noResults.remove();

  const isSent = payload.sender_id === myUserID;
  const div = document.createElement("div");
  div.className = `message ${isSent ? "sent" : "received"}`;
  div.innerHTML = `
    <span class="sender">${escapeHTML(payload.sender_name)}</span>
    ${payload.content ? `<p>${escapeHTML(payload.content)}</p>` : ""}
    ${messageImageHTML(payload.image_url)}
    <small>${formatDate(payload.created_at)}</small>
  `;
  chatWindow.appendChild(div);
  messageOffset++;
}

function scrollChatToBottom() {
  const chatWindow = document.getElementById("chatWindow");
  if (chatWindow) chatWindow.scrollTop = chatWindow.scrollHeight;
}

// =====================================================
// SCROLL PAGINATION WITH THROTTLE — Epic 8
// =====================================================
let lastScrollTime = 0;
const SCROLL_THROTTLE_MS = 500;

const SCROLL_TOP_THRESHOLD = 2;

function handleChatScroll(event) {
  const chatWindow = event.target;
  // scrollTop rarely lands on an exact 0 (fractional values on high-DPI
  // screens, brief negative overscroll on trackpads) — use a small threshold.
  if (chatWindow.scrollTop > SCROLL_TOP_THRESHOLD || !hasMoreMessages || isLoadingMessages) return;

  const now = Date.now();
  if (now - lastScrollTime < SCROLL_THROTTLE_MS) return;
  lastScrollTime = now;

  loadMessages(true);
}

// =====================================================
// IMAGE ATTACHMENT
// =====================================================
function showImagePreview(file) {
  const container = document.getElementById("imagePreviewContainer");
  const thumb = document.getElementById("imagePreviewThumb");
  if (!container || !thumb) return;

  thumb.src = URL.createObjectURL(file);
  container.classList.remove("hidden");
}

function clearImageSelection() {
  selectedImageFile = null;

  const imageInput = document.getElementById("messageImageInput");
  if (imageInput) imageInput.value = "";

  const container = document.getElementById("imagePreviewContainer");
  const thumb = document.getElementById("imagePreviewThumb");
  if (thumb && thumb.src) {
    URL.revokeObjectURL(thumb.src);
    thumb.src = "";
  }
  if (container) container.classList.add("hidden");
}

const attachImageButton = document.getElementById("attachImageButton");
const messageImageInput = document.getElementById("messageImageInput");
if (attachImageButton && messageImageInput) {
  attachImageButton.addEventListener("click", () => messageImageInput.click());
  messageImageInput.addEventListener("change", () => {
    const file = messageImageInput.files && messageImageInput.files[0];
    if (!file) return;
    selectedImageFile = file;
    showImagePreview(file);
  });
}

const removeImageButton = document.getElementById("removeImageButton");
if (removeImageButton) {
  removeImageButton.addEventListener("click", clearImageSelection);
}

async function uploadSelectedImage() {
  const formData = new FormData();
  formData.append("image", selectedImageFile);

  const response = await fetch("/api/messages/upload", {
    method: "POST",
    credentials: "include",
    body: formData,
  });

  if (!response.ok) {
    let errorMessage = "Could not upload image.";
    try {
      const error = await response.json();
      errorMessage = error.error || error.message || errorMessage;
    } catch (_) {}
    throw new Error(errorMessage);
  }

  const data = await response.json();
  return data.image_url;
}

// =====================================================
// SEND MESSAGE VIA WEBSOCKET — Epic 8
// =====================================================
async function sendMessage(event) {
  event.preventDefault();
  if (!currentChatUserID) {
    alert("Please select a conversation first.");
    return;
  }
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    alert("Not connected. Please wait and try again.");
    return;
  }

  const messageInput = document.getElementById("messageInput");
  const content = messageInput ? messageInput.value.trim() : "";
  if (!content && !selectedImageFile) return;

  let imageUrl = "";
  if (selectedImageFile) {
    try {
      imageUrl = await uploadSelectedImage();
    } catch (error) {
      console.error(error);
      alert(error.message || "Could not upload image.");
      return;
    }
  }

  socket.send(
    JSON.stringify({
      type: "chat_message",
      receiver_id: currentChatUserID,
      content: content,
      image_url: imageUrl,
    }),
  );

  messageInput.value = "";
  clearImageSelection();
}

const messageForm = document.getElementById("messageForm");
if (messageForm) {
  messageForm.addEventListener("submit", sendMessage);
}

// =====================================================
// NOTIFICATIONS — Epic 8
// =====================================================
let notificationCount = 0;

function showNotificationBadge(payload) {
  notificationCount++;

  const btn = document.getElementById("notificationButton");
  if (btn) {
    let badge = btn.querySelector(".notification-badge");
    if (!badge) {
      badge = document.createElement("span");
      badge.className = "notification-badge";
      badge.setAttribute("aria-label", "New notifications");
      btn.appendChild(badge);
    }
    badge.textContent = notificationCount;
    badge.classList.remove("hidden");
  }

  const notificationList = document.getElementById("notificationList");
  if (notificationList) {
    const li = document.createElement("li");
    li.innerHTML = `
      <strong>${escapeHTML(payload.sender_name)}</strong> sent you a message.
      <small>${formatDate(new Date().toISOString())}</small>
    `;
    notificationList.prepend(li);
  }
}

function clearNotificationBadge() {
  notificationCount = 0;
  const badge = document.querySelector(".notification-badge");
  if (badge) badge.classList.add("hidden");
}

const notificationButton = document.getElementById("notificationButton");
if (notificationButton) {
  notificationButton.addEventListener("click", () => {
    const panel = document.getElementById("notificationPanel");
    if (panel) {
      const isHidden = panel.classList.contains("hidden");
      panel.classList.toggle("hidden");
      notificationButton.setAttribute(
        "aria-expanded",
        isHidden ? "true" : "false",
      );
      if (isHidden) clearNotificationBadge();
    }
  });
}

// =====================================================
// TOAST — Epic 8
// =====================================================
function showToast(message) {
  const toast = document.getElementById("toast");
  if (!toast) return;
  toast.textContent = message;
  toast.classList.remove("hidden");
  toast.classList.add("show");
  setTimeout(() => {
    toast.classList.remove("show");
    toast.classList.add("hidden");
  }, 3000);
}

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
  loadUsers,
  openConversation,
};