import {renderCarousel, hideCarousel, renderTimeline} from './render.js'

let mainErr = document.getElementById('main-err')
let carousel = document.getElementById('carousel')

let posts = []
fetch('/api/all').then(resp => {
  if (!resp.ok) {
    console.log('failed to fetch posts', `status: ${resp.status}`)
    mainErr.innerText = "Oops, it borked. Come back later! 🛠️"
    return
  }

  resp.json().then(postsJson => {
    posts = postsJson.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    renderTimeline(posts)
  })
})

document.addEventListener('click', (e) => {
  // Show the carousel when clicking on pics
  const postToZoom = e.target.dataset.zoom
  if (postToZoom) {
    const post = posts.filter((p) => {return p.id === postToZoom})[0]
    renderCarousel(carousel, post)
    return
  }

  // Hide the carousel when 'clicking away'
  if (e.target.id !== 'carousel' && !e.target.classList.contains('large-pic')) {
    hideCarousel()
    return
  }
})


/* SECRET NAV HANDLERS */

let titleEl = document.getElementById('title')
let longTapTimer

function redirectToNewPost() {
  document.location = 'new-post.html'
}

titleEl.addEventListener('pointerdown', (e) => {
  e.preventDefault()
  longTapTimer = setTimeout(redirectToNewPost, 300)
})

function cancelTimer() {
  clearTimeout(longTapTimer)
}

titleEl.addEventListener('pointerup', cancelTimer)
titleEl.addEventListener('pointercancel', cancelTimer)
titleEl.addEventListener('pointermove', cancelTimer)
