
const API = 'http://localhost:8081'

function showToast(msg) {
    const t = document.createElement('div')
    t.textContent = msg
    t.className = 'fixed bottom-6 left-1/2 -translate-x-1/2 bg-black text-white text-xs px-4 py-2 rounded shadow z-50'
    document.body.appendChild(t)
    setTimeout(() => t.remove(), 2000)
}

// --- Interactions post card — l'API est seule source de vérité ---
function initPostCard(card, post) {
    let liked = !!post.liked_by_me
    let count = post.like_count || 0

    const likeBtn    = card.querySelector('.post-like-btn')
    const heart      = card.querySelector('.post-heart')
    const likesEl    = card.querySelector('.post-likes')
    const shareBtn   = card.querySelector('.post-share-btn')
    const commentBtn = card.querySelector('.post-comment-btn')

    const paint = () => {
        if (likesEl) likesEl.textContent = count
        heart?.setAttribute('fill', liked ? '#ef4444' : 'none')
        heart?.setAttribute('stroke', liked ? '#ef4444' : 'currentColor')
        likeBtn?.classList.toggle('text-red-500', liked)
    }
    paint()

    likeBtn?.addEventListener('click', (e) => {
        e.stopPropagation()
        liked = !liked
        count += liked ? 1 : -1
        paint()
        fetch(`${API}/api/posts/${post.id}/like`, {
            method: liked ? 'POST' : 'DELETE',
            credentials: 'include'
        }).catch(() => {
            liked = !liked
            count += liked ? 1 : -1
            paint()
        })
    })

    shareBtn?.addEventListener('click', (e) => {
        e.stopPropagation()
        navigator.clipboard.writeText(`${location.origin}/post?id=${post.id}`)
            .then(() => showToast('Lien copié !'))
    })

    commentBtn?.addEventListener('click', (e) => {
        e.stopPropagation()
        window.location.href = `/post?id=${post.id}`
    })
}

// Utilisateur courant en mémoire
let currentUser = null

// --- Session ---
async function checkAuth() {
    try {
        const res = await fetch(`${API}/api/auth/me`, { credentials: 'include' })
        if (!res.ok) return
        currentUser = await res.json()
        updateHeaderAuth()
    } catch (_) {}
}

function updateHeaderAuth() {
    if (!currentUser) return
    // Le header est géré côté serveur (Go template .User)
    // On affiche seulement le bouton "créer un post" si connecté
    const createBtn = document.getElementById('createPostBtn')
    if (createBtn) createBtn.style.display = 'inline-flex'
}

function initMenuBurger() {
    const menu = document.getElementById('menu-burger')
    const btn = document.querySelector('[data-burger]')
    if (!menu || !btn) return

    btn.addEventListener('click', (e) => {
        e.stopPropagation()
        menu.classList.toggle('hidden')
    })

    document.addEventListener('click', (e) => {
        if (!menu.contains(e.target)) menu.classList.add('hidden')
    })

    const logoutBtn = document.getElementById('menu-logout-btn')
    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => {
            fetch(`${API}/api/auth/logout`, { method: 'POST', credentials: 'include' })
                .finally(() => { window.location.href = '/login' })
        })
    }
}


function openEditProfileModal(profile = {}) {
    const modal = document.getElementById('modal-edit-profile')
    if (!modal) return

    document.getElementById('modal-pseudo').value = profile.pseudo || ''
    document.getElementById('modal-ville').value = profile.ville || ''
    document.getElementById('modal-bio').value = profile.bio || ''
    modal.classList.remove('hidden')
}

function closeEditProfileModal() {
    const modal = document.getElementById('modal-edit-profile')
    if (modal) modal.classList.add('hidden')
}

function initEditProfileModal() {
    const modal = document.getElementById('modal-edit-profile')
    if (!modal) return

    document.getElementById('modal-close-btn')?.addEventListener('click', closeEditProfileModal)

    modal.addEventListener('click', (e) => {
        if (e.target === modal) closeEditProfileModal()
    })

    document.getElementById('modal-save-btn')?.addEventListener('click', () => {
        const pseudo = document.getElementById('modal-pseudo').value.trim()
        const ville  = document.getElementById('modal-ville').value.trim()
        const bio    = document.getElementById('modal-bio').value.trim()

        // TODO: appel API PUT /api/users/me quand le back sera prêt
        console.log('Enregistrer:', { pseudo, ville, bio })

        const pseudoEl = document.getElementById('profile-pseudo')
        const villeEl  = document.querySelector('#profile-ville span:last-child')
        if (pseudoEl && pseudo) pseudoEl.textContent = pseudo
        if (villeEl  && ville)  villeEl.textContent  = ville

        closeEditProfileModal()
    })
}

function initContactsModal() {
    const modal = document.getElementById('modal-contacts')
    if (!modal) return

    document.getElementById('modal-contacts-close')?.addEventListener('click', () => modal.classList.add('hidden'))
    modal.addEventListener('click', (e) => { if (e.target === modal) modal.classList.add('hidden') })

    document.querySelector('[data-open-contacts]')?.addEventListener('click', (e) => {
        e.preventDefault()
        document.getElementById('menu-burger').classList.add('hidden')
        loadContacts()
        modal.classList.remove('hidden')
    })
}

function loadContacts(query = '') {
    fetch('/data/contacts.json')
        .then(r => r.json())
        .then(contacts => renderContacts(contacts, query))
}

function renderContacts(contacts, query = '') {
    const list = document.getElementById('contacts-list')
    if (!list) return
    const filtered = query
        ? contacts.filter(c => c.pseudo.toLowerCase().includes(query.toLowerCase()))
        : contacts

    list.innerHTML = ''
    filtered.forEach(contact => {
        const tpl = document.getElementById('contact-item').content.cloneNode(true)
        tpl.querySelector('.contact-pseudo').textContent = contact.pseudo

        const followBtn   = tpl.querySelector('.contact-follow-btn')
        const unfollowBtn = tpl.querySelector('.contact-unfollow-btn')

        followBtn.addEventListener('click', () => {
            contact.following = true
            followBtn.classList.add('opacity-50')
        })
        unfollowBtn.addEventListener('click', () => {
            contact.following = false
            unfollowBtn.classList.add('opacity-50')
        })

        list.appendChild(tpl)
    })

    document.getElementById('contacts-search')?.addEventListener('input', (e) => {
        renderContacts(contacts, e.target.value)
    }, { once: true })
}

function initFilterModal() {
    const modal = document.getElementById('modal-filter')
    if (!modal) return

    document.getElementById('modal-filter-close')?.addEventListener('click', () => modal.classList.add('hidden'))
    modal.addEventListener('click', (e) => { if (e.target === modal) modal.classList.add('hidden') })


    document.querySelector('[data-open-filter]')?.addEventListener('click', () => {
        modal.classList.remove('hidden')
    })

    document.getElementById('modal-filter-apply')?.addEventListener('click', () => {
        const sorts = [...document.querySelectorAll('input[name=sort]:checked')].map(i => i.value)
        const tags  = [...document.querySelectorAll('input[name=tag]:checked')].map(i => i.value)
        // TODO: appliquer le filtre sur les posts quand l'API sera prête
        console.log('Filtre:', { sorts, tags })
        modal.classList.add('hidden')
    })
}

// --- Dark / Light mode ---
function initTheme() {
    const isDark = localStorage.getItem('theme') === 'dark'
    if (isDark) applyDark()

    document.getElementById('theme-toggle')?.addEventListener('click', () => {
        const dark = document.documentElement.classList.contains('dark')
        dark ? applyLight() : applyDark()
    })
}

function applyDark() {
    document.documentElement.classList.add('dark')
    localStorage.setItem('theme', 'dark')
    document.getElementById('theme-label').textContent = 'LIGHT MODE'
    document.getElementById('theme-icon-moon')?.classList.add('hidden')
    document.getElementById('theme-icon-sun')?.classList.remove('hidden')
    document.getElementById('theme-pill')?.classList.replace('bg-gray-200', 'bg-black')
    document.getElementById('theme-dot')?.classList.add('translate-x-5')
}

function applyLight() {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('theme', 'light')
    document.getElementById('theme-label').textContent = 'DARK MODE'
    document.getElementById('theme-icon-moon')?.classList.remove('hidden')
    document.getElementById('theme-icon-sun')?.classList.add('hidden')
    document.getElementById('theme-pill')?.classList.replace('bg-black', 'bg-gray-200')
    document.getElementById('theme-dot')?.classList.remove('translate-x-5')
}

// --- Création de post ---
function initCreatePost() {
    const titre    = document.getElementById('cp-titre')
    const contenu  = document.getElementById('cp-contenu')
    const tagInput = document.getElementById('cp-tag-input')
    const tagsEl   = document.getElementById('cp-tags')
    const submit   = document.getElementById('cp-submit')
    const cancel   = document.getElementById('cp-cancel')
    const error    = document.getElementById('cp-error')
    const imageInput = document.getElementById('cp-image')
    const imageLabel = document.getElementById('cp-image-label')
    if (!submit) return

    imageInput?.addEventListener('change', () => {
        const file = imageInput.files[0]
        if (imageLabel) imageLabel.textContent = file ? file.name : 'Ajouter une image'
    })

    const tags = []

    function addTag(val) {
        val = val.trim().replace(/^#/, '').toLowerCase()
        if (!val || tags.includes(val)) return
        tags.push(val)
        const pill = document.createElement('span')
        pill.className = 'flex items-center gap-1 bg-black text-white px-2 py-0.5'
        pill.innerHTML = `#${val} <button data-tag="${val}" class="hover:opacity-70">×</button>`
        pill.querySelector('button').addEventListener('click', () => {
            tags.splice(tags.indexOf(val), 1)
            pill.remove()
            // Re-active le preset si c'en était un
            document.querySelector(`.cp-tag-preset[data-tag="${val}"]`)
                ?.classList.remove('bg-black', 'text-white', 'border-black')
        })
        tagsEl.appendChild(pill)
    }

    // Toggle tags prédéfinis
    document.querySelectorAll('.cp-tag-preset').forEach(preset => {
        preset.addEventListener('click', () => {
            const val = preset.dataset.tag
            if (tags.includes(val)) {
                tags.splice(tags.indexOf(val), 1)
                tagsEl.querySelector(`button[data-tag="${val}"]`)?.parentElement.remove()
                preset.classList.remove('bg-black', 'text-white', 'border-black')
            } else {
                addTag(val)
                preset.classList.add('bg-black', 'text-white', 'border-black')
            }
        })
    })

    // Tag personnalisé avec Entrée
    tagInput?.addEventListener('keydown', (e) => {
        if (e.key !== 'Enter') return
        e.preventDefault()
        addTag(tagInput.value)
        tagInput.value = ''
    })

    function resetForm() {
        titre.value = ''
        contenu.value = ''
        tagInput.value = ''
        tagsEl.innerHTML = ''
        tags.length = 0
        error.classList.add('hidden')
        document.querySelectorAll('.cp-tag-preset').forEach(p =>
            p.classList.remove('bg-black', 'text-white', 'border-black'))
    }

    cancel?.addEventListener('click', resetForm)

    submit.addEventListener('click', async () => {
        error.classList.add('hidden')
        const t = titre.value.trim()
        const c = contenu.value.trim()
        if (!t || !c) {
            error.textContent = 'Titre et description requis'
            error.classList.remove('hidden')
            return
        }
        try {
            const res = await fetch(`${API}/api/posts`, {
                method: 'POST',
                credentials: 'include',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ titre: t, contenu: c, tags })
            })
            if (!res.ok) {
                const data = await res.json()
                throw new Error(data.error || 'Erreur')
            }
            // Reset et recharge les posts
            const post = await res.json()
            // Upload image si sélectionnée
            if (imageInput?.files[0]) {
                const fd = new FormData()
                fd.append('image', imageInput.files[0])
                await fetch(`${API}/api/posts/${post.id}/image`, {
                    method: 'POST',
                    credentials: 'include',
                    body: fd
                })
            }
            resetForm()
            showToast('Post publié !')
            if (typeof reloadPosts === 'function') reloadPosts()
        } catch (err) {
            error.textContent = err.message
            error.classList.remove('hidden')
        }
    })
}

document.addEventListener('DOMContentLoaded', () => {
    initTheme()
    initMenuBurger()
    initEditProfileModal()
    initContactsModal()
    initFilterModal()
    initCreatePost()
    checkAuth()
    initPasswordToggle()
})

function initPasswordToggle() {
    const toggleButtons = document.querySelectorAll('.toggle-password')
    toggleButtons.forEach(btn => {
        btn.addEventListener('click', (e) => {
            e.preventDefault()
            const targetId = btn.getAttribute('data-target')
            const input = document.getElementById(targetId)
            if (!input) return

            const isPassword = input.type === 'password'
            input.type = isPassword ? 'text' : 'password'

            // Changer l'icône (ajouter une classe ou changer la couleur)
            btn.classList.toggle('text-gray-700')
            btn.classList.toggle('text-gray-500')
        })
    })
}
