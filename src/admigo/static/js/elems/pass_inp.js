;"use strict";
class PassInp {
  constructor() {
    window.fun.ready(this.handler.bind(this));
  }

  shpa_icon_click(e) {
    e.preventDefault();
    e.stopPropagation();

    const btn = e.currentTarget;
    const par = window.fun.parent(btn, '.shpa-container');
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
      if (window.fun.once(btn, 'shpa_icon_click')) return;
      btn.addEventListener('click', this.shpa_icon_click.bind(this));
    });
  }
}

new PassInp()
