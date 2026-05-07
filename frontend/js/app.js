
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

        // Mise à jour locale de l'affichage
        const pseudoEl = document.getElementById('profile-pseudo')
        const villeEl  = document.querySelector('#profile-ville span:last-child')
        if (pseudoEl && pseudo) pseudoEl.textContent = pseudo
        if (villeEl  && ville)  villeEl.textContent  = ville

        closeEditProfileModal()
    })
}

document.addEventListener('DOMContentLoaded', () => {
    initMenuBurger()
    initEditProfileModal()
})
