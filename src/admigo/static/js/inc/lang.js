;"use strict";

class Lang {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));
    this.lame = null;
  }

  handler() {
    let stat = document.documentElement.getAttribute("stat");
    let lang = document.documentElement.getAttribute("lang");
    let tra = document.documentElement.getAttribute("tra");

    if (parseInt(tra) == 0) {
      return;
    }

    fetch('/' + stat + '/tra/' + lang + '/messages.json')
      .then((re) => {
        if (!re.ok) return null;
        return re.json()
      })
      .then((new_json) => {
        this.lame = new_json;
      });
  }

  re(key) {
    if (!this.lame) {
      return key;
    }

    let ret = key;

    ret = this.lame[key];

    if (!ret) {
      return key;
    }

    return ret;
  }
}

window.lang = new Lang();
