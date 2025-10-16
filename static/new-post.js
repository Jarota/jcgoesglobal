import {renderCarousel, hideCarousel, renderPost} from './render.js'

let captionInput = document.getElementById('caption')
let dateInput = document.getElementById('post-date')
let fileInput = document.getElementById('image-upload')
let preview = document.getElementById('preview')
let carousel = document.getElementById('carousel')

// obj to preview
let post = {
  id:         'preview-id',
  author:     'Me',
  caption:    captionInput.value,
  pics:       [],
  date:       new Date().toISOString()
}

captionInput.addEventListener('change', updatePreviewCaption)
function updatePreviewCaption() {
  post.caption = captionInput.value
  preview.innerHTML = '' // replace previous preview
  preview.appendChild(renderPost(post))
}

dateInput.addEventListener('change', updatePreviewDate)
function updatePreviewDate() {
  post.date = dateInput.value
  preview.innerHTML = '' // replace previous preview
  preview.appendChild(renderPost(post))
}

fileInput.addEventListener('change', updatePreviewPics)
function updatePreviewPics() {
  let loaded = 0
  const total = fileInput.files.length
  for (const file of fileInput.files) {
    const reader = new FileReader()

    reader.onload = function (e) {
      post.pics.push(e.target.result)

      loaded++
      if (loaded === total) {
        preview.innerHTML = '' // replace previous preview
        preview.appendChild(renderPost(post))
      }
    }
    reader.readAsDataURL(file)
  }
}

document.addEventListener('click', (e) => {
  // Show the carousel when clicking on pics
  const postToZoom = e.target.dataset.zoom
  if (postToZoom) {
    renderCarousel(carousel, post)
    return
  }

  // Hide the carousel when 'clicking away'
  if (e.target.id !== 'carousel' && !e.target.classList.contains('large-pic')) {
    hideCarousel()
    return
  }
})

function initRender() {
  dateInput.value = post.date.split('T')[0];
  updatePreviewCaption()
  updatePreviewDate()
}
initRender()
