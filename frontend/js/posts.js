function timeAgo(dateStr) {
    const diff = Math.floor((Date.now() - new Date(dateStr)) / 60000)
    if (diff < 60) return `il y a ${diff} minute${diff > 1 ? 's' : ''}`
    const h = Math.floor(diff / 60)
    if (h < 24) return `il y a ${h} heure${h > 1 ? 's' : ''}`
    const d = Math.floor(h / 24)
    return `il y a ${d} jour${d > 1 ? 's' : ''}`
}

function createPostCard(post) {
    const card = document.getElementById('post-card').content.cloneNode(true)
    const authorEl = card.querySelector('.post-author')
    authorEl.textContent = post.author_pseudo
    authorEl.href = `/profile?id=${post.user_id}`
    card.querySelector('.post-title').textContent = post.titre
    card.querySelector('.post-content').textContent = post.contenu
    card.querySelector('.post-date').textContent = timeAgo(post.date_publication)
    card.querySelector('.post-tags').textContent = (post.tags || []).map(t => `#${t}`).join(' ')
    card.querySelector('.post-likes').textContent = post.like_count

    if (post.image_url) {
        const img = document.createElement('img')
        img.src = `${window.location.origin}${post.image_url}`
        img.className = 'max-w-full max-h-96 object-contain object-left mt-1'
        img.alt = post.titre
        const titleEl = card.querySelector('.post-title')
        titleEl?.insertAdjacentElement('afterend', img)
    }

    const article = card.querySelector('article')
    article.style.cursor = 'pointer'
    article.addEventListener('click', (e) => {
        if (e.target.closest('button')) return
        window.location.href = `/post?id=${post.id}`
    })
    initPostCard(article, post)
    return card
}

function createTopPostCard(post) {
    const card = document.getElementById('top-post-card').content.cloneNode(true)
    const a = card.querySelector('.top-post-author')
    a.textContent = post.author_pseudo
    a.href = `/profile?id=${post.user_id}`
    card.querySelector('.top-post-content').textContent = post.contenu
    card.querySelector('.top-post-date').textContent = timeAgo(post.date_publication)
    return card
}

const postsFilter = { q: '', sort: '', tags: [] }

function loadPosts() {
    const container = document.getElementById('posts-container')
    if (container) container.innerHTML = ''

    const params = new URLSearchParams()
    if (postsFilter.q)          params.set('q', postsFilter.q)
    if (postsFilter.sort)       params.set('sort', postsFilter.sort)
    if (postsFilter.tags.length) params.set('tags', postsFilter.tags.join(','))

    const url = `${API}/api/posts?${params.toString()}`
    console.log('loadPosts - URL de la requête:', url)
    console.log('loadPosts - postsFilter:', postsFilter)

    fetch(url, { credentials: 'include' })
        .then(r => {
            console.log('loadPosts - Réponse reçue:', r.status)
            return r.json()
        })
        .then(posts => {
            console.log('loadPosts - Posts reçus:', posts)
            if (!Array.isArray(posts)) {
                console.error('loadPosts - Réponse non-array:', posts)
                return
            }
            if (posts.length === 0 && container) {
                container.innerHTML = '<p class="text-sm text-gray-400 text-center py-8">Aucun post trouvé</p>'
                console.log('loadPosts - Aucun post trouvé')
                return
            }
            console.log('loadPosts - Ajout de', posts.length, 'posts')
            posts.forEach(post => container?.appendChild(createPostCard(post)))
        })
        .catch(err => {
            console.error('loadPosts - Erreur fetch:', err)
        })
}

function loadTopPosts() {
    const topContainer = document.getElementById('top-posts-container')
    if (topContainer) topContainer.innerHTML = ''

    fetch(`${API}/api/posts/top?limit=6`, { credentials: 'include' })
        .then(r => r.json())
        .then(posts => {
            if (!Array.isArray(posts)) return
            posts.forEach(post => topContainer?.appendChild(createTopPostCard(post)))
        })
}

function reloadPosts() { loadPosts(); loadTopPosts() }

function renderUserSuggestions(users) {
    const box = document.getElementById('search-suggestions')
    if (!box) return
    if (!users || users.length === 0) { box.classList.add('hidden'); return }
    box.innerHTML = users.map(u => `
        <a href="/profile?id=${u.id}" class="flex items-center gap-3 px-4 py-2 hover:bg-gray-100 border-b border-gray-100 last:border-0">
            <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 shrink-0 overflow-hidden">
                ${u.photo_url
                    ? `<img src="${window.location.origin}${u.photo_url}" class="w-full h-full object-cover">`
                    : `<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>`
                }
            </div>
            <div class="flex flex-col">
                <span class="text-sm font-bold">${u.pseudo}</span>
                <span class="text-xs text-gray-400">${u.nom || ''}</span>
            </div>
        </a>
    `).join('')
    box.classList.remove('hidden')
}

function searchUsers(q) {
    if (!q) { renderUserSuggestions([]); return }
    fetch(`${API}/api/search?q=${encodeURIComponent(q)}&limit=10`, { credentials: 'include' })
        .then(r => r.json())
        .then(results => {
            const users = (results || []).filter(r => r.type === 'user').map(r => r.user)
            renderUserSuggestions(users)
        })
}

function initSearch() {
    const input = document.getElementById('search-input')
    const box   = document.getElementById('search-suggestions')
    if (!input) return
    let t
    input.addEventListener('input', () => {
        clearTimeout(t)
        t = setTimeout(() => {
            const q = input.value.trim()
            postsFilter.q = q
            loadPosts()
            searchUsers(q)
        }, 300)
    })
    input.addEventListener('focus', () => {
        if (input.value.trim()) searchUsers(input.value.trim())
    })
    document.addEventListener('click', (e) => {
        if (box && !box.contains(e.target) && e.target !== input) box.classList.add('hidden')
    })
}

if (typeof checkAuth === 'function') {
    checkAuth().then(() => { 
        loadPosts()
        loadTopPosts()
    })
} else {
    loadPosts()
    loadTopPosts()
}

// Initialiser la recherche une fois que le DOM est prêt
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initSearch)
} else {
    initSearch()
}
