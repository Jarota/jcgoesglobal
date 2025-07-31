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
  // Show the carousel when clicking a zoom button
  const postToZoom = e.target.dataset.zoom
  if (postToZoom) {
    renderCarousel(postToZoom)
    return
  }

  // Hide the carousel when 'clicking away'
  if (e.target.id !== 'carousel' && !e.target.classList.contains('large-pic')) {
    carousel.classList.add('hidden')
    carousel.classList.remove('flex')
    document.body.style.overflowY = '' // Re-enable scrolling too
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
  carousel.classList.add('flex')
  carousel.classList.remove('hidden')

  // Also prevent scrolling in the background
  document.body.style.overflowY = 'hidden'
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
        <p class="caption">${post.caption}</p>
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
    pics.forEach((pic) => {
      picsHtml += `<img class="pic" src="assets/pics/${pic}">`
    })

    return `
      <div class="pics flex">
        ${picsHtml}
      </div>
      <button class="zoom-btn" data-zoom="${id}">
        <i class="fa-solid fa-up-right-and-down-left-from-center" data-zoom="${id}"></i>
      </button>
    `
}
