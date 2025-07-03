import { posts } from './data.js'

let carousel = document.getElementById('carousel')

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

  document.getElementById('timeline').innerHTML = timelineHtml
}

function renderPics(id, pics) {
    if (pics.length === 0) {
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

renderTimeline()
