let elements = {};
let nodes = new Map();
window._GOUI_ELEMENTS = elements;
let randInt = () => Math.floor(Math.random() * 2e9) + 1;
let generateId = () => {
    let id = randInt();
    while (elements[id]) id = randInt();
    return id;
};
let createElement = (tag, clicks) => {
    let id = generateId();
    let el = document.createElement(tag);
    elements[id] = el;
    if (clicks) nodes.set(el, id);
    return id;
};
let createElementNS = (tag, ns, clicks) => {
    let id = generateId();
    let el = document.createElementNS(ns, tag);
    elements[id] = el;
    if (clicks) nodes.set(el, id);
    return id;
};
let createTextNode = text => {
    let id = generateId();
    elements[id] = document.createTextNode(text);
    return id;
};
let decoder = new TextDecoder();
let memory;
let exports;
let getString = (addr, len) => len ? decoder.decode(memory.buffer.slice(addr, addr + len)) : '';
let go = new Go();
Object.assign(go.importObject.gojs, {
    createElement: (addr, len, clicks) => createElement(getString(addr, len), clicks),
    createTd: clicks => createElement('td', clicks),
    createTr: clicks => createElement('tr', clicks),
    createSpan: clicks => createElement('span', clicks),
    createA: clicks => createElement('a', clicks),
    createDiv: clicks => createElement('div', clicks),
    createTable: clicks => createElement('table', clicks),
    createTbody: clicks => createElement('tbody', clicks),
    createH1: clicks => createElement('h1', clicks),
    createButton: clicks => createElement('button', clicks),
    createElementNS: (addr, len, addr2, len2, clicks) => createElementNS(getString(addr, len), getString(addr2, len2), clicks),
    createTextNode: (addr, len) => createTextNode(getString(addr, len)),
    appendChild: (parent, child) => {
        elements[parent].appendChild(elements[child]);
    },
    setStr: (node, addr, len, addr2, len2) => {
        elements[node][getString(addr, len)] = getString(addr2, len2);
    },
    setTextContent: (node, addr, len) => {
        elements[node].textContent = getString(addr, len);
    },
    setData: (node, addr, len) => {
        elements[node].data = getString(addr, len);
    },
    setClass: (node, addr, len) => {
        elements[node].className = getString(addr, len);
    },
    setAriaHidden: (node, bool) => {
        elements[node].ariaHidden = !!bool;
    },
    setBool: (node, addr, len, bool) => {
        elements[node][getString(addr, len)] = !!bool;
    },
    replaceWith: (oldNode, newNode) => {
        elements[oldNode].replaceWith(elements[newNode]);
        delete elements[oldNode];
    },
    removeAttribute: (node, addr, len) => {
        elements[node].removeAttribute(getString(addr, len));
    },
    removeNode: node => {
        if (elements[node]) {
            
            elements[node].remove();
        } else {
            console.log('rm', node);

        }
    },
    disposeNode: node => {
        let el = elements[node];
        if (nodes.has(el)) nodes.delete(el);
        delete elements[node];
    },
    cloneNode: node => {
        let id = generateId();
        elements[id] = elements[node].cloneNode(true);
        return id;
    },
    moveBefore: (parent, nextKeyMatch, start, movingDomNode) => {
        let mdm = elements[movingDomNode];
        let par = elements[parent];
        let curr = par.childNodes[start];
        if (mdm === curr) return;
        let oldPos = mdm.nextSibling;
        par.insertBefore(mdm, curr);
        if (curr !== par.lastChild && !nextKeyMatch) {
            par.insertBefore(curr, oldPos);
        }
    },
    mount: (node, addr, len) => {
        let root = document.querySelector(getString(addr, len));
        root.appendChild(elements[node]);
        root.addEventListener('click', e => {
            window._GOUI_EVENT = e;
            let target = e.target;
            while (target && target != root) {
                let node = nodes.get(target);
                if (node) {
                    exports.callClickListener(node);
                }
                target = target.parentNode;
            }
        });
    },
});

WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject).then(o => {
    let instance = o.instance;
    exports = instance.exports;
    memory = exports.memory;
    go.run(instance);
});
