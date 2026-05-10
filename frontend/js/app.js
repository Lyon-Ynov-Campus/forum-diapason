// JS commun à toutes les pages
// Utilitaires 

async function apiPost(url, data) {
    const res = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    })
    const json = await res.json()
    if (!res.ok) throw new Error(json.error || 'Erreur serveur')
    return json
}

async function apiDelete(url) {
    const res = await fetch(url, { method: 'DELETE' })
    const json = await res.json()
    if (!res.ok) throw new Error(json.error || 'Erreur serveur')
    return json
}

function showError(formEl, message) {
    let el = formEl.querySelector('.form-error')
    if (!el) {
        el = document.createElement('p')
        el.className = 'form-error text-red-500 text-sm mt-1'
        formEl.prepend(el)
    }
    el.textContent = '⚠ ' + message
}

// Auth - vérifie si connecté au chargement de la page
let currentUser = null

async function checkAuth() {
    try {
        const res = await fetch('/api/auth/me')
        if (res.ok) {
            currentUser = await res.json()
            showUserNav()
        } else {
            showGuestNav()
        }
    } catch {
        showGuestNav()
    }
}

function showUserNav() {
    document.getElementById('nav-guest')?.classList.add('hidden')
    document.getElementById('nav-user')?.classList.remove('hidden')
}

function showGuestNav() {
    document.getElementById('nav-guest')?.classList.remove('hidden')
    document.getElementById('nav-user')?.classList.add('hidden')
}

//Register

const registerForm = document.getElementById('register-form')
if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault()
        const btn = registerForm.querySelector('button[type=submit]')
        btn.disabled = true
        btn.textContent = 'Chargement...'

        try {
            await apiPost('/api/auth/register', {
                nom:      registerForm.querySelector('[name=nom]').value,
                pseudo:   registerForm.querySelector('[name=pseudo]').value,
                email:    registerForm.querySelector('[name=email]').value,
                password: registerForm.querySelector('[name=password]').value,
            })
            window.location.href = '/'
        } catch (err) {
            showError(registerForm, err.message)
            btn.disabled = false
            btn.textContent = 'Register'
        }
    })
}
//  menu dropdown
const menuBtn      = document.getElementById('user-menu-btn')
const menuDropdown = document.getElementById('user-menu-dropdown')

if (menuBtn && menuDropdown) {
    menuBtn.addEventListener('click', (e) => {
        e.stopPropagation()
        menuDropdown.classList.toggle('hidden')
    })
    // Fermer en cliquant ailleurs
    document.addEventListener('click', () => {
        menuDropdown.classList.add('hidden')
    })
}


//Login

const loginForm = document.getElementById('login-form')
if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault()
        const btn = loginForm.querySelector('button[type=submit]')
        btn.disabled = true
        btn.textContent = 'Chargement...'

        try {
            await apiPost('/api/auth/login', {
                email_or_pseudo: loginForm.querySelector('[name=email_or_pseudo]').value,
                password:        loginForm.querySelector('[name=password]').value,
            })
            window.location.href = '/'
        } catch (err) {
            showError(loginForm, err.message)
            btn.disabled = false
            btn.textContent = 'Sign in'
        }
    })
}
//Create Post 

const createPostBtn = document.getElementById('create-post-btn')
if (createPostBtn) {
    const tags = []

    // Gestion des tags
    const tagInput  = document.getElementById('tag-input')
    const tagList   = document.getElementById('tag-list')

    if (tagInput) {
        tagInput.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' || e.key === ',') {
                e.preventDefault()
                const val = tagInput.value.trim().replace(/^#/, '')
                if (val && !tags.includes(val)) {
                    tags.push(val)
                    renderTags()
                }
                tagInput.value = ''
            }
        })
    }

    function renderTags() {
        if (!tagList) return
        tagList.innerHTML = tags.map((t, i) => `
            <span class="flex items-center gap-1 border border-gray-300 px-2 py-0.5">
                #${t}
                <button type="button" data-i="${i}">×</button>
            </span>`).join('')
        tagList.querySelectorAll('button').forEach(btn => {
            btn.addEventListener('click', () => {
                tags.splice(+btn.dataset.i, 1)
                renderTags()
            })
        })
    }

    createPostBtn.addEventListener('click', async () => {
        const titre   = document.getElementById('post-titre')?.value.trim()
        const contenu = document.getElementById('post-contenu')?.value.trim()
        if (!titre || !contenu) return

        try {
            await apiPost('/api/posts', { titre, contenu, tags })
            window.location.reload()
        } catch (err) {
            alert(err.message)
        }
    })
}

//  create post
const postTags = []

const tagInput = document.getElementById('tag-input')
const tagList  = document.getElementById('tag-list')

if (tagInput) {
    tagInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' || e.key === ',') {
            e.preventDefault()
            const val = tagInput.value.trim().replace(/^#/, '').toLowerCase()
            if (val && !postTags.includes(val)) {
                postTags.push(val)
                renderPostTags()
            }
            tagInput.value = ''
        }
    })
}

function renderPostTags() {
    if (!tagList) return
    tagList.innerHTML = postTags.map((t, i) => `
        <span class="flex items-center gap-1 border border-gray-300 px-2 py-0.5 text-xs">
            #${t}
            <button type="button" data-i="${i}" class="hover:text-red-500">×</button>
        </span>`).join('')
    tagList.querySelectorAll('button').forEach(btn => {
        btn.addEventListener('click', () => {
            postTags.splice(+btn.dataset.i, 1)
            renderPostTags()
        })
    })
}

document.getElementById('create-post-btn')?.addEventListener('click', async () => {
    if (!currentUser) {
        window.location.href = '/login'
        return
    }

    const titre   = document.getElementById('post-titre')?.value.trim()
    const contenu = document.getElementById('post-contenu')?.value.trim()

    if (!titre)   { alert('Le titre est obligatoire'); return }
    if (!contenu) { alert('Le contenu est obligatoire'); return }

    try {
        const post = await apiPost('/api/posts', { titre, contenu, tags: postTags })

        // Réinitialiser le formulaire
        document.getElementById('post-titre').value   = ''
        document.getElementById('post-contenu').value = ''
        postTags.length = 0
        renderPostTags()

        // Recharger le feed
        loadPosts()
    } catch (err) {
        alert(err.message)
    }
})
//  Feed — charger les posts depuis l'API

function timeAgo(dateStr) {
    const diff = Math.floor((Date.now() - new Date(dateStr)) / 60000)
    if (diff < 1)  return 'à l\'instant'
    if (diff < 60) return `il y a ${diff} minute${diff > 1 ? 's' : ''}`
    const h = Math.floor(diff / 60)
    if (h < 24) return `il y a ${h} heure${h > 1 ? 's' : ''}`
    const d = Math.floor(h / 24)
    return `il y a ${d} jour${d > 1 ? 's' : ''}`
}

function createPostCard(post) {
    const tpl = document.getElementById('post-card')
    if (!tpl) return null
    const card = tpl.content.cloneNode(true)

    card.querySelector('.post-author').textContent  = post.author_pseudo || post.author || '?'
    card.querySelector('.post-title').textContent   = post.titre  || post.title   || ''
    card.querySelector('.post-content').textContent = post.contenu || post.content || ''
    card.querySelector('.post-date').textContent    = timeAgo(post.date_publication || post.created_at)
    card.querySelector('.post-tags').textContent    = (post.tags || []).map(t => `#${t}`).join(' ')

    // Like button
    const likeBtn   = card.querySelector('.like-btn')
    const likeCount = card.querySelector('.post-likes')
    if (likeBtn && likeCount) {
        likeCount.textContent = post.like_count ?? post.likes ?? 0
        if (post.liked_by_me) likeBtn.classList.add('text-red-500')

        likeBtn.addEventListener('click', async () => {
            if (!currentUser) { window.location.href = '/login'; return }
            try {
                if (post.liked_by_me) {
                    await apiDelete(`/api/posts/${post.id}/like`)
                    post.liked_by_me = false
                    post.like_count  = (post.like_count || 1) - 1
                    likeBtn.classList.remove('text-red-500')
                } else {
                    await apiPost(`/api/posts/${post.id}/like`, {})
                    post.liked_by_me = true
                    post.like_count  = (post.like_count || 0) + 1
                    likeBtn.classList.add('text-red-500')
                }
                likeCount.textContent = post.like_count
            } catch (err) {
                console.error(err.message)
            }
        })
    }

    return card
}

async function loadPosts() {
    const container = document.getElementById('posts-container')
    if (!container) return

    try {
        const res   = await fetch('/api/posts')
        const posts = await res.json()

        container.innerHTML = ''
        posts.forEach(post => {
            const card = createPostCard(post)
            if (card) container.appendChild(card)
        })

        // Top posts
        const topContainer = document.getElementById('top-posts-container')
        if (topContainer) {
            topContainer.innerHTML = ''
            ;[...posts]
                .sort((a, b) => (b.like_count || 0) - (a.like_count || 0))
                .slice(0, 6)
                .forEach(post => {
                    const tpl = document.getElementById('top-post-card')
                    if (!tpl) return
                    const card = tpl.content.cloneNode(true)
                    card.querySelector('.top-post-author').textContent  = post.author_pseudo || '?'
                    card.querySelector('.top-post-content').textContent = post.contenu || ''
                    card.querySelector('.top-post-date').textContent    = timeAgo(post.date_publication)
                    topContainer.appendChild(card)
                })
        }
    } catch (err) {
        console.error('Erreur chargement posts:', err)
    }
}
// init
checkAuth().then(() => {
    loadPosts()
})