async function getBase64FromBlob(blob) {
  const base64 = await new Promise((resolve) => {
    const reader = new FileReader();
    reader.readAsDataURL(blob);
    reader.onloadend = () => {
      const base64data = reader.result;
      resolve(base64data);
    }
  });
  return base64
}

async function getFilenameFromResponse(response) {
  const header = response.headers.get('Content-Disposition');
  const parts = header?.split(';');
  filename = parts[1].split('=')[1];
  return filename
}

async function getImageFromApi(url){
  const response = await fetch(url)
  const blob = await response.blob()
  const filename = await getFilenameFromResponse(response)
  console.log({blob,filename})
  if (filename.includes(".pdf")){
    return window.open(url)
  }
  const base64 = await getBase64FromBlob(blob)
  const image_window = window.open(url, "_blank")
  image_window.document.write(`<html><head></head><body><img width="300" src=${base64} alt="loading" /></body></html>`);
}
