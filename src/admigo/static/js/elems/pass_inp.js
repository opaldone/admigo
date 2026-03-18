;"use strict";
class PassInp {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));
  }

  shpa_icon_click(e) {
    e.preventDefault();
    e.stopPropagation();

    const btn = e.currentTarget;
    const par = this.fun.parent(btn, '.shpa-container');
    const is_show = par.classList.contains('sh');

    if (is_show) {
      par.classList.remove('sh')
      return false;
    }

    par.classList.add('sh');
    return false;
  }

  handler() {
    document.querySelectorAll('.shpa-icon').forEach(btn => {
      if (this.fun.once(btn, 'shpa_icon_click')) return;
      btn.addEventListener('click', this.shpa_icon_click.bind(this));
    });
  }
}

new PassInp()
