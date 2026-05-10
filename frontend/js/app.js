// JS commun à toutes les pages
// ── Utilitaires ──────────────────────────────────────────────

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

// ── Register ─────────────────────────────────────────────────

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

// ── Login ────────────────────────────────────────────────────

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