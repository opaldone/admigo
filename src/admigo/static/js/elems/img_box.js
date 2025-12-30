;"use strict";
class ImgBox {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));
  }

  handler() {
    document.querySelectorAll('.img-box-container').forEach(el => {
      let img_tag = el.querySelector('.img-tag');
      let file_input = el.querySelector('.img-selector');
      let start_crop = el.querySelector('.start-crop');
      let end_crop = el.querySelector('.end-crop');
      let size = el.querySelector('.crop-size');

      el.data = {
        'crop': null,
        'size': size
      };

      start_crop.data = {
        'img_tag': img_tag,
        'file_input': file_input,
        'par': el
      };

      end_crop.data = {
        'img_tag': img_tag,
        'file_input': file_input,
        'par': el
      };

      file_input.data = {
        'img_tag': img_tag,
        'par': el
      };

      file_input.addEventListener('change', this.file_input_change.bind(this));
      start_crop.addEventListener('click', this.start_crop_click.bind(this));
      end_crop.addEventListener('click', this.end_crop_click.bind(this));

      el.addEventListener('crop', this.el_crop.bind(this));
    });
  }

  el_crop(ev) {
    let el = ev.currentTarget;
    let sz = el.data.size;
    if (!sz) return;

    let st = Math.round(ev.detail.width) + '<sig>&centerdot;</sig>' + Math.round(ev.detail.height);

    sz.innerHTML = st;
  }

  file_input_change(ev) {
    let inp = ev.currentTarget;
    let fl = ev.target.files[0];
    let img_tag = inp.data.img_tag;
    let par = inp.data.par;

    if (!fl) return;

    inp.data.file_name = fl.name;
    inp.data.file_type = fl.type;

    const reader = new FileReader();

    reader.onload = (e) => {
      img_tag.src = e.target.result
      img_tag.classList.add('show');
      par.classList.add('ex-img');
    }

    reader.readAsDataURL(fl);
  }

  clearCrop(el) {
    let dt = el.data;

    dt.img_tag.classList.add('show');
    dt.par.data.crop.destroy();
    dt.par.data.crop = null;

    dt.file_input.classList.remove('hid');
    dt.par.classList.remove('cropping');
  }

  start_crop_click(ev) {
    let btn = ev.currentTarget;
    let dt = btn.data;
    let par = dt.par

    if (dt.img_tag.getAttribute("src").length == 0) {
      return;
    }

    if (par.data.crop != null) {
      this.clearCrop(btn)
      return;
    }

    dt.file_input.classList.add('hid');
    par.classList.add('cropping');

    let ratio = par.getAttribute('data-ratio');

    let co = {
      'viewMode': 1,
    };

    if (ratio) {
      co['aspectRatio'] = 1;
    }

    par.data.crop = new Cropper(dt.img_tag, co);
  }

  async img_type(src_in) {
    const res = await fetch(src_in);
    const bl = await res.blob();
    return bl.type;
  }

  async get_ft(dt) {
    let ret = {
      'file_name': 'img.png',
      'file_type': 'image/png'
    };
    let fn = dt.file_input.data.file_name;
    let ft = dt.file_input.data.file_type;

    if (fn && ft) {
      ret.file_name = fn;
      ret.file_type = ft;
      return ret;
    }

    let isrc = dt.img_tag.getAttribute("src")
    if (isrc.length == 0) {
      return ret;
    }

    fn = isrc.split("/").reverse()[0];
    ft = await this.img_type(isrc);

    ret.file_name = fn;
    ret.file_type = ft;
    return ret;
  }

  async end_crop_click(ev) {
    let btn = ev.currentTarget;
    let dt = btn.data;

    if (!dt.par.data.crop) return;

    let ft = await this.get_ft(dt);

    dt.par.data.crop.getCroppedCanvas().toBlob((bl) => {
      const re = new FileReader();
      re.onload = (e) => {
        dt.img_tag.src = e.target.result

        this.clearCrop(btn);
      }
      re.readAsDataURL(bl);

      let file = new File([bl], ft.file_name, {type: ft.file_type, lastModified:new Date().getTime()});
      let container = new DataTransfer();
      container.items.add(file);
      dt.file_input.files = container.files;
    }, ft.file_type);
  }
}

new ImgBox()
