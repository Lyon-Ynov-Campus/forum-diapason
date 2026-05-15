function timeAgo(dateStr) {
    const normalized = dateStr.replace(' +0000 UTC', 'Z').replace(' ', 'T')
    const diff = Math.floor((Date.now() - new Date(normalized)) / 60000)
    if (diff < 1) return "à l'instant"
    if (diff < 60) return `il y a ${diff} minute${diff > 1 ? 's' : ''}`
    const h = Math.floor(diff / 60)
    if (h < 24) return `il y a ${h} heure${h > 1 ? 's' : ''}`
    const d = Math.floor(h / 24)
    return `il y a ${d} jour${d > 1 ? 's' : ''}`
}

function createComment(comment) {
    const tpl = document.getElementById('comment-card').content.cloneNode(true)

    const authorEl = tpl.querySelector('.comment-author')
    authorEl.textContent = comment.author_pseudo
    authorEl.href = `/profile?pseudo=${encodeURIComponent(comment.author_pseudo)}`

    tpl.querySelector('.comment-content').textContent = comment.contenu
    tpl.querySelector('.comment-date').textContent = timeAgo(comment.date)

    let liked = false
    let count = 0
    const likesEl = tpl.querySelector('.comment-likes')
    const likeBtn = tpl.querySelector('.comment-like-btn')
    likesEl.textContent = count
    likeBtn?.addEventListener('click', () => {
        liked = !liked
        count += liked ? 1 : -1
        likesEl.textContent = count
        likeBtn.classList.toggle('text-red-500', liked)
        const svg = likeBtn.querySelector('svg')
        svg?.setAttribute('fill', liked ? '#ef4444' : 'none')
        svg?.setAttribute('stroke', liked ? '#ef4444' : 'currentColor')
    })

    const replyBtn    = tpl.querySelector('.comment-reply-btn')
    const replyForm   = tpl.querySelector('.comment-reply-form')
    const replySubmit = tpl.querySelector('.comment-reply-submit')
    const replyInput  = replyForm?.querySelector('input')
    const repliesEl   = tpl.querySelector('.comment-replies')

    replyBtn?.addEventListener('click', () => {
        replyForm?.classList.toggle('hidden')
        if (replyInput && !replyForm?.classList.contains('hidden')) {
            replyInput.value = `@${comment.author_pseudo} `
            replyInput.focus()
        }
    })

    replySubmit?.addEventListener('click', () => {
        const val = replyInput?.value.trim()
        if (!val) return
        fetch(`${API}/api/posts/${postId}/comments`, {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ contenu: val })
        })
        .then(r => r.json())
        .then(reply => {
            const card = document.getElementById('reply-card').content.cloneNode(true)
            const a = card.querySelector('.reply-author')
            a.textContent = reply.author_pseudo
            a.href = `/profile?pseudo=${encodeURIComponent(reply.author_pseudo)}`
            card.querySelector('.reply-content').textContent = reply.contenu
            card.querySelector('.reply-date').textContent = timeAgo(reply.date)
            repliesEl.appendChild(card)
            replyInput.value = ''
            replyForm.classList.add('hidden')
        })
    })

    return tpl
}

const postId = parseInt(new URLSearchParams(window.location.search).get('id'))

if (postId) {
    Promise.all([
        fetch(`${API}/api/posts/${postId}`, { credentials: 'include' }).then(r => r.json()),
        fetch(`${API}/api/posts/${postId}/comments`, { credentials: 'include' }).then(r => r.json())
    ]).then(([post, comments]) => {
        render(post)
        const list = document.getElementById('comments-list')
        ;(comments || []).forEach(c => list.appendChild(createComment(c)))
        setupNewComment(postId, list)
    }).catch(() => {
        document.getElementById('post-title').textContent = 'Post introuvable'
    })
}

function render(post) {
    const authorEl = document.getElementById('post-author')
    authorEl.textContent = post.author_pseudo
    authorEl.href = `/profile?id=${post.user_id}`

    document.getElementById('post-title').textContent   = post.titre
    document.getElementById('post-content').textContent = post.contenu
    document.getElementById('post-date').textContent    = timeAgo(post.date_publication)
    document.getElementById('post-tags').textContent    = (post.tags || []).map(t => `#${t}`).join(' ')
    document.getElementById('post-likes').textContent   = post.like_count

    // Afficher l'image si elle existe
    if (post.image_url) {
        const imageContainer = document.getElementById('post-image-container')
        if (imageContainer) {
            const img = document.createElement('img')
            img.src = `http://localhost:8080${post.image_url}`
            img.className = 'w-full h-full object-cover'
            img.alt = post.titre
            imageContainer.innerHTML = ''
            imageContainer.appendChild(img)
        }
    }

    // Like + share
    const likeBtn = document.querySelector('.post-detail-like')
    const heart   = document.querySelector('.post-detail-heart')
    let liked = post.liked_by_me
    let count = post.like_count

    if (liked) {
        heart?.setAttribute('fill', '#ef4444')
        heart?.setAttribute('stroke', '#ef4444')
        likeBtn?.classList.add('text-red-500')
    }

    likeBtn?.addEventListener('click', () => {
        liked = !liked
        count += liked ? 1 : -1
        document.getElementById('post-likes').textContent = count
        heart?.setAttribute('fill', liked ? '#ef4444' : 'none')
        heart?.setAttribute('stroke', liked ? '#ef4444' : 'currentColor')
        likeBtn.classList.toggle('text-red-500', liked)
        const method = liked ? 'POST' : 'DELETE'
        fetch(`${API}/api/posts/${post.id}/like`, { method, credentials: 'include' })
    })

    document.querySelector('.post-detail-share')?.addEventListener('click', () => {
        navigator.clipboard.writeText(location.href).then(() => showToast('Lien copié !'))
    })
}

function setupNewComment(postId, list) {
    document.getElementById('new-comment-submit')?.addEventListener('click', () => {
        const input = document.getElementById('new-comment-input')
        const val = input.value.trim()
        if (!val) return
        fetch(`${API}/api/posts/${postId}/comments`, {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ contenu: val })
        })
        .then(r => r.json())
        .then(comment => {
            list.appendChild(createComment(comment))
            input.value = ''
            list.scrollTop = list.scrollHeight
        })
    })

    document.getElementById('new-comment-input')?.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') document.getElementById('new-comment-submit')?.click()
    })
}
