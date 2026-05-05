const authTabLogin = document.getElementById('auth-tab-login');
const authTabRegister = document.getElementById('auth-tab-register');
const headerSigninBtn = document.getElementById('header-signin-btn');
const headerSignupBtn = document.getElementById('header-signup-btn');
const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const authMessage = document.getElementById('auth-message');

const storageKey = 'diapason-forum-state';

function loadState() {
  const raw = localStorage.getItem(storageKey);
  if (!raw) {
    return { currentUser: null, users: [], posts: [] };
  }

  try {
    return JSON.parse(raw);
  } catch (error) {
    return { currentUser: null, users: [], posts: [] };
  }
}

function saveState(state) {
  localStorage.setItem(storageKey, JSON.stringify(state));
}

function createNotification(message, type = 'info') {
  if (!authMessage) return;
  authMessage.textContent = message;
  authMessage.classList.remove('hidden');
  authMessage.dataset.state = type;
}

function clearNotification() {
  if (!authMessage) return;
  authMessage.textContent = '';
  authMessage.classList.add('hidden');
}

function setAuthTab(tab) {
  authTabLogin.classList.toggle('active', tab === 'login');
  authTabRegister.classList.toggle('active', tab === 'register');

  document.getElementById('login-panel').classList.toggle('hidden', tab !== 'login');
  document.getElementById('register-panel').classList.toggle('hidden', tab !== 'register');
  clearNotification();
}

function loginUser(event) {
  event.preventDefault();
  const state = loadState();
  const email = document.getElementById('login-email').value.trim();
  const password = document.getElementById('login-password').value.trim();

  if (!email || !password) {
    createNotification('Merci de remplir tous les champs.', 'error');
    return;
  }

  const user = state.users.find(
    (item) => item.email.toLowerCase() === email.toLowerCase() && item.password === password
  );

  if (!user) {
    createNotification('Email ou mot de passe incorrect.', 'error');
    return;
  }

  state.currentUser = { email: user.email, pseudo: user.pseudo, name: user.name };
  saveState(state);
  window.location.href = '/';
}

function registerUser(event) {
  event.preventDefault();
  const state = loadState();
  const email = document.getElementById('register-email').value.trim();
  const pseudo = document.getElementById('register-pseudo').value.trim();
  const name = document.getElementById('register-name').value.trim();
  const password = document.getElementById('register-password').value.trim();
  const passwordConfirm = document.getElementById('register-password-confirm').value.trim();

  if (!email || !pseudo || !name || !password || !passwordConfirm) {
    createNotification('Merci de remplir tous les champs.', 'error');
    return;
  }

  if (password !== passwordConfirm) {
    createNotification('Les mots de passe ne correspondent pas.', 'error');
    return;
  }

  if (state.users.some((item) => item.email.toLowerCase() === email.toLowerCase())) {
    createNotification('Cet email est déjà utilisé.', 'error');
    return;
  }

  state.users.push({ email, pseudo, name, password });
  state.currentUser = { email, pseudo, name };
  saveState(state);
  window.location.href = '/';
}

authTabLogin.addEventListener('click', () => setAuthTab('login'));
authTabRegister.addEventListener('click', () => setAuthTab('register'));
headerSigninBtn.addEventListener('click', () => setAuthTab('login'));
headerSignupBtn.addEventListener('click', () => setAuthTab('register'));
loginForm.addEventListener('submit', loginUser);
registerForm.addEventListener('submit', registerUser);

// Check URL parameter for initial tab
const params = new URLSearchParams(window.location.search);
const initialTab = params.get('tab') || 'login';
setAuthTab(initialTab);
