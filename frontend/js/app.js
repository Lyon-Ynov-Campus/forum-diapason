const storageKey = 'diapason-forum-state';

const ui = {
  authTabLogin: document.getElementById('auth-tab-login'),
  authTabRegister: document.getElementById('auth-tab-register'),
  authFormLogin: document.getElementById('login-form'),
  authFormRegister: document.getElementById('register-form'),
  authMessage: document.getElementById('auth-message'),
  createPostMessage: document.getElementById('create-post-message'),
  createPostForm: document.getElementById('create-post-form'),
  postFeed: document.getElementById('post-feed'),
  topPosts: document.getElementById('top-posts'),
  createPostPanel: document.getElementById('create-post-panel'),
  authContainer: document.getElementById('auth-container'),
  signOutButton: document.getElementById('sign-out-button'),
  searchInput: document.getElementById('search-input'),
  welcomeMessage: document.getElementById('welcome-message'),
  authButton: document.getElementById('auth-button'),
};

const defaultPosts = [
  {
    id: 1,
    author: 'Anna',
    title: 'Nouvelle pédale d\'effet : mes impressions',
    content: 'J\'ai testé une nouvelle pédale de delay et elle change vraiment le son. Qui d\'autre l\'a essayée ?',
    tags: ['#guitar', '#studio'],
    likes: 24,
    createdAt: 'Il y a 15 minutes',
    likedByMe: false,
  },
  {
    id: 2,
    author: 'Tom',
    title: 'Recherche d\'un sample de basse chaude',
    content: 'Je cherche un sample de basse analogique pour un morceau lo-fi. Des pistes ?',
    tags: ['#bass', '#sample'],
    likes: 12,
    createdAt: 'Il y a 35 minutes',
    likedByMe: false,
  },
  {
    id: 3,
    author: 'Leila',
    title: 'Comment préparer une session mix ?',
    content: 'Je voudrais optimiser mon workflow de mix. Quels sont vos presets et astuces ?',
    tags: ['#mix', '#workflow'],
    likes: 18,
    createdAt: 'Il y a 1 heure',
    likedByMe: false,
  },
  {
    id: 4,
    author: 'Marcus',
    title: 'Les meilleurs DAW en 2026',
    content: 'Après 10 ans de production, voici mon comparatif des meilleures stations de travail audio numériques. Ableton Live reste mon préféré pour la flexibilité.',
    tags: ['#daw', '#production', '#logiciel'],
    likes: 42,
    createdAt: 'Il y a 2 heures',
    likedByMe: false,
  },
  {
    id: 5,
    author: 'Sophie',
    title: 'Synth synthé synthétique : où trouver les meilleurs presets ?',
    content: 'J\'aimerais recommandations de sites ou communautés pour télécharger des presets de qualité pour Serum et Sylenth1.',
    tags: ['#synth', '#presets', '#serum'],
    likes: 9,
    createdAt: 'Il y a 3 heures',
    likedByMe: false,
  },
  {
    id: 6,
    author: 'Jules',
    title: 'Mon premier EP est en ligne !',
    content: 'Après 6 mois de travail, j\'ai enfin terminé mon premier EP de musique électronique. C\'est disponible sur Spotify et Bandcamp. Merci à la communauté pour les retours !',
    tags: ['#ep', '#électronique', '#release'],
    likes: 56,
    createdAt: 'Il y a 5 heures',
    likedByMe: false,
  },
  {
    id: 7,
    author: 'Raphaël',
    title: 'Podcast musique : recommandations ?',
    content: 'Je cherche des podcasts intéressants sur la production musicale et l\'industrie musicale en général. Vos suggestions ?',
    tags: ['#podcast', '#musique', '#audio'],
    likes: 14,
    createdAt: 'Il y a 6 heures',
    likedByMe: false,
  },
  {
    id: 8,
    author: 'Nina',
    title: 'Critique du nouvel album de SolidState',
    content: 'Le nouvel album est franchement incroyable. Les arrangements sont sophistiqués et la production est impeccable. Je n\'ai aucun regret d\'avoir acheté la version vinyle.',
    tags: ['#album', '#critique', '#vinyle'],
    likes: 31,
    createdAt: 'Il y a 7 heures',
    likedByMe: false,
  },
  {
    id: 9,
    author: 'David',
    title: 'Setup home studio 2026 : mon avis',
    content: 'Après avoir investi dans un nouveau setup, je veux partager mon configuration avec vous. Audiobox, Neumann U87, moniteurs Adam A7X et traitement acoustique.',
    tags: ['#homestudio', '#setup', '#équipement'],
    likes: 67,
    createdAt: 'Il y a 8 heures',
    likedByMe: false,
  },
  {
    id: 10,
    author: 'Emma',
    title: 'Techniques de mastering : les basiques',
    content: 'Un guide pour les débutants en mastering. Voici les 5 étapes essentielles que j\'ai apprises en travaillant avec des ingénieurs professionnels.',
    tags: ['#mastering', '#tutoriel', '#audio'],
    likes: 43,
    createdAt: 'Il y a 10 heures',
    likedByMe: false,
  },
  {
    id: 11,
    author: 'Lucas',
    title: 'Faire de la musique gratuitement : guide complet',
    content: 'Vous voulez créer de la musique mais vous n\'avez pas de budget ? Voici les meilleurs logiciels gratuits et open-source pour débuter.',
    tags: ['#gratuit', '#logiciel', '#débutant'],
    likes: 89,
    createdAt: 'Il y a 12 heures',
    likedByMe: false,
  },
  {
    id: 12,
    author: 'Clara',
    title: 'Synthèse FM : tutoriel avancé',
    content: 'La synthèse FM peut être intimidante. Je propose un cours détaillé sur comment maîtriser cette technique puissante et créer des sons uniques.',
    tags: ['#fm', '#synthèse', '#tutoriel'],
    likes: 37,
    createdAt: 'Il y a 14 heures',
    likedByMe: false,
  },
  {
    id: 13,
    author: 'Pierre',
    title: 'Collaboration musicale : comment trouver des partenaires ?',
    content: 'Je cherche d\'autres musiciens intéressés par la collaboration. Je fais principalement du synthwave et du darkwave. Des intéressés ?',
    tags: ['#collaboration', '#synthwave', '#darkwave'],
    likes: 22,
    createdAt: 'Il y a 16 heures',
    likedByMe: false,
  },
  {
    id: 14,
    author: 'Victoria',
    title: 'Réduction du bruit : mes astuces',
    content: 'Pour enregistrer à la maison sans studio acoustique, voici comment j\'ai réduit considérablement le bruit ambiant avec peu de moyens.',
    tags: ['#recording', '#acoustique', '#astuces'],
    likes: 28,
    createdAt: 'Il y a 18 heures',
    likedByMe: false,
  },
  {
    id: 15,
    author: 'Olivier',
    title: 'Les meilleurs plugins VST gratuits',
    content: 'Compilation des VST gratuits les meilleurs du marché. J\'utilise ces plugins régulièrement et ils rivalisent avec les versions payantes.',
    tags: ['#vst', '#plugins', '#gratuit'],
    likes: 76,
    createdAt: 'Il y a 20 heures',
    likedByMe: false,
  },
];

function loadState() {
  const raw = localStorage.getItem(storageKey);
  if (!raw) {
    return {
      currentUser: null,
      users: [],
      posts: defaultPosts,
    };
  }

  try {
    return JSON.parse(raw);
  } catch (error) {
    return {
      currentUser: null,
      users: [],
      posts: defaultPosts,
    };
  }
}

function saveState(state) {
  localStorage.setItem(storageKey, JSON.stringify(state));
}

function createNotification(container, message, type = 'info') {
  if (!container) return;
  container.textContent = message;
  container.classList.remove('hidden');
  container.dataset.state = type;
}

function clearNotification(container) {
  if (!container) return;
  container.textContent = '';
  container.classList.add('hidden');
}

function formatTags(raw) {
  return raw
    .split(',')
    .map((tag) => tag.trim())
    .filter(Boolean)
    .map((tag) => (tag.startsWith('#') ? tag : `#${tag}`));
}

function renderTopPosts(posts) {
  ui.topPosts.innerHTML = '';
  const top = [...posts].sort((a, b) => b.likes - a.likes).slice(0, 5);

  top.forEach((post) => {
    const item = document.createElement('div');
    item.className = 'top-post';
    item.innerHTML = `
      <div class="top-avatar">${post.author[0] || 'U'}</div>
      <div class="top-post-content">
        <p class="top-post-title">${post.title}</p>
        <p class="top-post-text">${post.content.slice(0, 72)}${post.content.length > 72 ? '...' : ''}</p>
        <p class="inline-label">${post.likes} likes · ${post.createdAt}</p>
      </div>
    `;
    ui.topPosts.appendChild(item);
  });
}

function renderPosts(posts, query = '') {
  ui.postFeed.innerHTML = '';
  const normalizedQuery = query.trim().toLowerCase();

  const visible = posts.filter((post) => {
    if (!normalizedQuery) return true;
    const content = `${post.title} ${post.content} ${post.tags.join(' ')}`.toLowerCase();
    return content.includes(normalizedQuery);
  });

  if (visible.length === 0) {
    const empty = document.createElement('div');
    empty.className = 'notification';
    empty.textContent = 'Aucune publication ne correspond à cette recherche.';
    ui.postFeed.appendChild(empty);
    return;
  }

  visible.forEach((post) => {
    const article = document.createElement('article');
    article.className = 'post-card';
    article.innerHTML = `
      <div class="post-head">
        <div class="author">
          <div class="avatar">${post.author[0] || 'U'}</div>
          <div>
            <div>${post.author}</div>
            <div class="post-meta">${post.createdAt}</div>
          </div>
        </div>
        <div class="post-meta">${post.tags.join(' ')}</div>
      </div>
      <div class="card-body">
        <h2 class="post-title">${post.title}</h2>
        <p class="post-text">${post.content}</p>
      </div>
      <div class="post-footer">
        <button class="like-button ${post.likedByMe ? 'liked' : ''}" data-id="${post.id}">
          ${post.likedByMe ? '♥' : '♡'} ${post.likes}
        </button>
      </div>
    `;
    ui.postFeed.appendChild(article);
  });

  ui.postFeed.querySelectorAll('.like-button').forEach((button) => {
    button.addEventListener('click', () => {
      const postId = Number(button.dataset.id);
      toggleLike(postId);
    });
  });
}

function toggleLike(postId) {
  const state = loadState();
  const post = state.posts.find((item) => item.id === postId);
  if (!post) return;
  post.likedByMe = !post.likedByMe;
  post.likes += post.likedByMe ? 1 : -1;
  saveState(state);
  renderPosts(state.posts, ui.searchInput.value);
  renderTopPosts(state.posts);
}

function updateAuthUI() {
  const state = loadState();
  const isLoggedIn = Boolean(state.currentUser);

  ui.createPostPanel.classList.toggle('hidden', !isLoggedIn);
  ui.authContainer.classList.toggle('hidden', isLoggedIn);
  ui.signOutButton.classList.toggle('hidden', !isLoggedIn);
  ui.authButton.classList.toggle('hidden', isLoggedIn);
  ui.welcomeMessage.textContent = isLoggedIn ? `Bienvenue, ${state.currentUser.pseudo} !` : '';
}

function setAuthTab(tab) {
  ui.authTabLogin.classList.toggle('active', tab === 'login');
  ui.authTabRegister.classList.toggle('active', tab === 'register');
  document.getElementById('login-panel').classList.toggle('hidden', tab !== 'login');
  document.getElementById('register-panel').classList.toggle('hidden', tab !== 'register');
  clearNotification(ui.authMessage);
}

function registerUser(event) {
  event.preventDefault();
  const state = loadState();
  const email = document.getElementById('register-email').value.trim();
  const pseudo = document.getElementById('register-pseudo').value.trim();
  const password = document.getElementById('register-password').value.trim();

  if (!email || !pseudo || !password) {
    createNotification(ui.authMessage, 'Merci de remplir tous les champs.', 'error');
    return;
  }

  if (state.users.some((user) => user.email.toLowerCase() === email.toLowerCase())) {
    createNotification(ui.authMessage, 'Cet email est déjà utilisé.', 'error');
    return;
  }

  state.users.push({ email, pseudo, password });
  state.currentUser = { email, pseudo };
  saveState(state);
  ui.authFormRegister.reset();
  updateAuthUI();
  renderPosts(state.posts, ui.searchInput.value);
  createNotification(ui.authMessage, 'Inscription réussie ! Vous êtes connecté.', 'success');
}

function loginUser(event) {
  event.preventDefault();
  const state = loadState();
  const email = document.getElementById('login-email').value.trim();
  const password = document.getElementById('login-password').value.trim();

  const user = state.users.find(
    (item) => item.email.toLowerCase() === email.toLowerCase() && item.password === password
  );

  if (!user) {
    createNotification(ui.authMessage, 'Email ou mot de passe incorrect.', 'error');
    return;
  }

  state.currentUser = { email: user.email, pseudo: user.pseudo };
  saveState(state);
  ui.authFormLogin.reset();
  updateAuthUI();
  createNotification(ui.authMessage, 'Connexion réussie !', 'success');
}

function signOut() {
  const state = loadState();
  state.currentUser = null;
  saveState(state);
  updateAuthUI();
  createNotification(ui.authMessage, 'Vous êtes déconnecté.', 'info');
}

function createPost(event) {
  event.preventDefault();
  const state = loadState();
  if (!state.currentUser) {
    createNotification(ui.createPostMessage, 'Connectez-vous pour poster.', 'error');
    return;
  }

  const title = document.getElementById('post-title').value.trim();
  const content = document.getElementById('post-content').value.trim();
  const tags = formatTags(document.getElementById('post-tags').value);

  if (!title || !content) {
    createNotification(ui.createPostMessage, 'Titre et description sont requis.', 'error');
    return;
  }

  const nextId = state.posts.length ? Math.max(...state.posts.map((post) => post.id)) + 1 : 1;
  state.posts.unshift({
    id: nextId,
    author: state.currentUser.pseudo,
    title,
    content,
    tags,
    likes: 0,
    createdAt: 'À l\'instant',
    likedByMe: false,
  });

  saveState(state);
  ui.createPostForm.reset();
  renderPosts(state.posts, ui.searchInput.value);
  renderTopPosts(state.posts);
  createNotification(ui.createPostMessage, 'Publication ajoutée avec succès !', 'success');
}

function init() {
  setAuthTab('login');
  ui.authTabLogin.addEventListener('click', () => setAuthTab('login'));
  ui.authTabRegister.addEventListener('click', () => setAuthTab('register'));
  ui.authFormRegister.addEventListener('submit', registerUser);
  ui.authFormLogin.addEventListener('submit', loginUser);
  ui.createPostForm.addEventListener('submit', createPost);
  ui.signOutButton.addEventListener('click', signOut);
  ui.authButton.addEventListener('click', () => {
    ui.authContainer.scrollIntoView({ behavior: 'smooth' });
  });
  ui.searchInput.addEventListener('input', () => {
    const state = loadState();
    renderPosts(state.posts, ui.searchInput.value);
  });

  updateAuthUI();
  const state = loadState();
  renderPosts(state.posts, ui.searchInput.value);
  renderTopPosts(state.posts);
}

init();
