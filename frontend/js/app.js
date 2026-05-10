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

function showError(formEl, message) {
    let el = formEl.querySelector('.form-error')
    if (!el) {
        el = document.createElement('p')
        el.className = 'form-error text-red-500 text-sm'
        formEl.prepend(el)
    }
    el.textContent = message
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