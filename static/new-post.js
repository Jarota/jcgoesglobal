import {renderCarousel, hideCarousel, renderPost} from './render.js'

let preview = document.getElementById('preview')
// obj to preview
let post = {
  id:        'preview-id',
  author:    '',
  caption:   '',
  pics:      [],
  createdAt: ''
}

let captionInput = document.getElementById('caption')
captionInput.addEventListener('change', updatePreviewCaption)
function updatePreviewCaption() {
  post.caption = captionInput.value
  preview.innerHTML = renderPost(post)
}

let fileInput = document.getElementById('image-upload')
fileInput.addEventListener('change', updatePreviewPics)
function updatePreviewPics() {
  for (const file of fileInput.files) {
    post.pics.push(file.name)
  }
  preview.innerHTML = renderPost(post)
}

let carousel = document.getElementById('carousel')

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
