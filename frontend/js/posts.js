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
        img.src = `http://localhost:8080${post.image_url}`
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

function createUserCard(user) {
    const div = document.createElement('div')
    div.className = 'border border-gray-200 p-4 hover:bg-gray-50 cursor-pointer transition'
    div.style.cursor = 'pointer'
    div.innerHTML = `
        <div class="flex items-center gap-3">
            <img src="${user.photo_url || 'https://via.placeholder.com/40'}" alt="${user.pseudo}" class="w-10 h-10 rounded-full object-cover">
            <div>
                <p class="font-bold text-sm">${user.pseudo}</p>
                <p class="text-xs text-gray-500">👤 Utilisateur</p>
            </div>
        </div>
    `
    div.addEventListener('click', () => {
        window.location.href = `/profile?id=${user.id}`
    })
    return div
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

    // Si y a une recherche textuelle ET pas de tags/tri -> utiliser l'API de recherche globale
    if (postsFilter.q && !postsFilter.sort && postsFilter.tags.length === 0) {
        return loadSearchResults()
    }

    // Sinon utiliser l'API posts avec filtres
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

function loadSearchResults() {
    const container = document.getElementById('posts-container')
    const url = `${API}/api/search?q=${encodeURIComponent(postsFilter.q)}`
    console.log('loadSearchResults - URL:', url)

    fetch(url, { credentials: 'include' })
        .then(r => r.json())
        .then(results => {
            console.log('loadSearchResults - Résultats:', results)
            if (!Array.isArray(results)) return
            
            if (results.length === 0 && container) {
                container.innerHTML = '<p class="text-sm text-gray-400 text-center py-8">Aucun résultat trouvé pour "<strong>' + postsFilter.q + '</strong>"</p>'
                return
            }
            
            results.forEach(result => {
                if (result.type === 'user' && result.user) {
                    container?.appendChild(createUserCard(result.user))
                } else if (result.type === 'post' && result.post) {
                    container?.appendChild(createPostCard(result.post))
                }
            })
        })
        .catch(err => {
            console.error('loadSearchResults - Erreur fetch:', err)
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

function initSearch() {
    const input = document.getElementById('search-input')
    if (!input) return
    let t
    input.addEventListener('input', () => {
        clearTimeout(t)
        t = setTimeout(() => {
            postsFilter.q = input.value.trim()
            loadPosts()
        }, 300)
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
