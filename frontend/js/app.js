
const API = 'http://localhost:8081'

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
    const signIn   = document.querySelector('a[href="/login"]')
    const register = document.querySelector('a[href="/register"]')
    if (signIn)   signIn.style.display   = 'none'
    if (register) register.style.display = 'none'

    // Affiche le pseudo dans la nav
    const nav = document.querySelector('nav.flex')
    if (nav && !document.getElementById('nav-user')) {
        const span = document.createElement('a')
        span.id        = 'nav-user'
        span.href      = `/profile?id=${currentUser.id}`
        span.textContent = currentUser.pseudo
        span.className = 'text-sm font-bold hover:underline'
        nav.prepend(span)
    }
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

document.addEventListener('DOMContentLoaded', () => {
    initTheme()
    initMenuBurger()
    initEditProfileModal()
    initContactsModal()
    initFilterModal()
    checkAuth()
})
