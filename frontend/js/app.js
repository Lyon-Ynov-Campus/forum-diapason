
function showToast(msg) {
    const t = document.createElement('div')
    t.textContent = msg
    t.className = 'fixed bottom-6 left-1/2 -translate-x-1/2 bg-black text-white text-xs px-4 py-2 rounded shadow z-50'
    document.body.appendChild(t)
    setTimeout(() => t.remove(), 2000)
}


function initPostCard(card, post) {
    let liked = false
    let count = post.likes

    const likeBtn    = card.querySelector('.post-like-btn')
    const heart      = card.querySelector('.post-heart')
    const likesEl    = card.querySelector('.post-likes')
    const shareBtn   = card.querySelector('.post-share-btn')
    const commentBtn = card.querySelector('.post-comment-btn')

    likeBtn?.addEventListener('click', (e) => {
        e.stopPropagation()
        liked = !liked
        count += liked ? 1 : -1
        likesEl.textContent = count
        heart.setAttribute('fill', liked ? '#ef4444' : 'none')
        heart.setAttribute('stroke', liked ? '#ef4444' : 'currentColor')
        likeBtn.classList.toggle('text-red-500', liked)
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

// --- Menu burger ---
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
            document.cookie = 'session_id=; Max-Age=-1; path=/'
            window.location.href = '/login'
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


function initAuthForms() {
    document.querySelectorAll('input[type=password]').forEach(input => {
        const btn = document.createElement('button')
        btn.type = 'button'
        btn.textContent = '👁'
        btn.className = 'absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 text-sm'
        input.parentElement.style.position = 'relative'
        input.style.paddingRight = '2rem'
        input.parentElement.appendChild(btn)
        btn.addEventListener('click', () => {
            input.type = input.type === 'password' ? 'text' : 'password'
        })
    })

    // Validation register : mots de passe identiques
    const pw  = document.querySelector('input[name=password]')
    const pw2 = document.querySelector('input[name=password_confirm]')
    const submitBtn = document.querySelector('form button[type=submit]')
    if (pw && pw2 && submitBtn) {
        const check = () => {
            const mismatch = pw2.value && pw.value !== pw2.value
            pw2.style.borderColor = mismatch ? '#ef4444' : ''
            submitBtn.disabled = !!mismatch
        }
        pw.addEventListener('input', check)
        pw2.addEventListener('input', check)
    }
}

document.addEventListener('DOMContentLoaded', () => {
    initTheme()
    initMenuBurger()
    initEditProfileModal()
    initContactsModal()
    initFilterModal()
    initAuthForms()
})
