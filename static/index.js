let carousel = document.getElementById('carousel')

let posts = []
fetch('/api/all').then(resp => {
  if (!resp.ok) {
    console.log('failed to fetch posts', `status: ${resp.status}`)
    return
  }

  resp.json().then(postsJson => {
    posts = postsJson.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    renderTimeline()
  })
})

document.addEventListener('click', (e) => {
  // Show the carousel when clicking on pics
  const postToZoom = e.target.dataset.zoom
  if (postToZoom) {
    renderCarousel(postToZoom)
    return
  }

  // Hide the carousel when 'clicking away'
  if (e.target.id !== 'carousel' && !e.target.classList.contains('large-pic')) {
    hideCarousel()
    return
  }
})

function renderCarousel(id) {
  const post = posts.filter((p) => {return p.id === id})[0]

  let zoomedHtml = ''
  post.pics.forEach((path) => {
    zoomedHtml += `
      <img src="assets/pics/${path}" class="large-pic">
    `
  })
  carousel.innerHTML = zoomedHtml
  carousel.classList.remove('hidden')
  carousel.classList.add('shown')

  // Also prevent scrolling in the background
  document.body.style.overflowY = 'hidden'
}

function hideCarousel() {
    carousel.classList.remove('shown')
    carousel.classList.add('hidden')
    document.body.style.overflowY = '' // Re-enable scrolling too
}

function renderTimeline() {
  let timelineHtml = ''

  let firstPost = true
  posts.forEach((post) => {
    let verticalLine = '<div class="vr"></div>'
    if (firstPost) {
      verticalLine = ''
      firstPost = false
    }

    timelineHtml += `
      ${verticalLine}
      <div class="post" id="${post.id}">
        <div class="caption-cntr">
          <p class="caption">${post.caption}</p>
        </div>
        ${renderPics(post.id, post.pics)}
      </div>
    `
  })

  if (timelineHtml === '') {
    timelineHtml = `<p>No posts (yet)!</p>`
  }

  document.getElementById('timeline').innerHTML = timelineHtml
}

function renderPics(id, pics) {
    if (!pics || pics.length === 0) {
      return ''
    }

    let picsHtml = ''
    let degrees = -10
    let z = -1
    pics.forEach((pic) => {
      picsHtml += `<img class="pic"
        src="assets/pics/${pic}" data-zoom="${id}"
        style="transform: rotate(${degrees}deg); z-index: ${z};"
      >`
      degrees = (degrees > 0 ? degrees + 10 : degrees - 10) * -1
      z -= 1
    })

    return `
      <div class="pics" data-zoom="${id}">
        ${picsHtml}
      </div>
    `
}

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
