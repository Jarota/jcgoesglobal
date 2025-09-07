function renderCarousel(carousel, post) {
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
  document.body.style.overflow = 'hidden'
  document.documentElement.style.overflow = 'hidden'
}

function hideCarousel() {
  carousel.classList.remove('shown')
  carousel.classList.add('hidden')
  document.body.style.overflow = '' // Re-enable scrolling too
  document.documentElement.style.overflow = ''
}

function renderTimeline(posts) {
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
      ${renderPost(post)}
    `
  })

  if (timelineHtml === '') {
    timelineHtml = `<p>No posts (yet)!</p>`
  }

  document.getElementById('timeline').innerHTML = timelineHtml
}

function renderPost(post) {
  return `<div class="post" id="${post.id}">
      <div class="caption-cntr">
        <p class="caption">${post.caption}</p>
      </div>
      ${renderPics(post.id, post.pics)}
    </div>
  `  
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

export {renderCarousel, hideCarousel, renderTimeline, renderPost}
