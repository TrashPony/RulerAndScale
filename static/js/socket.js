let ws;

function Connect() {
    ws = new WebSocket("ws://" + window.location.host + "/ws");
    console.log("Websocket - status: " + ws.readyState);

    ws.onopen = function () {
        setInterval(function () {
            ws.send(JSON.stringify({event: "Debug"}));
        }, 100);
        console.log("Connection opened..." + this.readyState);
    };

    ws.onmessage = function (msg) {
        DrawTop(JSON.parse(msg.data));
        DrawFront(JSON.parse(msg.data));
    };

    ws.onerror = function (msg) {
        console.log("Error occured sending..." + msg.data);
    };

    ws.onclose = function (msg) {
        console.log("Disconnected lobby - status " + this.readyState);
    };
}

let Scale = 4;

function DrawTop(data) {
    console.log(data)
    let canvas = document.getElementById("canvasTop");
    canvas.width = 648;
    canvas.height = 400;

    let ctx = canvas.getContext("2d");
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // рисуем полоску обозначающая размер СМ

    ctx.font = "8px serif";
    ctx.fillText("1 см", Scale + Scale / 2, 10);

    ctx.beginPath();
    ctx.moveTo(Scale, 15);
    ctx.lineTo(2 * Scale, 15);
    ctx.stroke();

    // рисуем область измерения
    ctx.strokeStyle = "blue";
    ctx.rect(
        canvas.width / 2 - (data.ruler_option.width_max * (Scale / 2)),
        10,
        data.ruler_option.width_max * Scale,
        data.ruler_option.length_max * Scale,
    );
    ctx.stroke();

    // рисуем платформу весов
    ctx.strokeStyle = "black";
    ctx.fillStyle = "rgba(153,161,163,0.71)";
    ctx.roundRect(
        canvas.width / 2 - (data.scale_platform.width * (Scale / 2)),
        10,
        data.scale_platform.width * Scale,
        data.scale_platform.height * Scale,
        5);
    ctx.stroke();
    ctx.fill();

    // рисуем левый датчик
    ctx.fillStyle = "rgb(0,255,0)";
    if (data.indication.left < 0) ctx.fillStyle = "rgb(255,0,0)";
    let left = {x: canvas.width / 2 - (data.ruler_option.width_max * (Scale / 2)) - 7, y: 17};
    ctx.fillRect(left.x, left.y, 7, 7);

    // рисуем правый датчик
    ctx.fillStyle = "rgb(0,255,0)";
    if (data.indication.right < 0) ctx.fillStyle = "rgb(255,0,0)";
    let right = {x: canvas.width / 2 + (data.ruler_option.width_max * (Scale / 2)), y: 20};
    ctx.fillRect(right.x, right.y, 7, 7);

    // рисуем нижний датчи
    ctx.fillStyle = "rgb(0,255,0)";
    if (data.indication.back < 0) ctx.fillStyle = "rgb(255,0,0)";
    let back = {x: canvas.width / 2 - 3.5, y: (data.ruler_option.length_max * (Scale / 2)) * 2 + 10};
    ctx.fillRect(back.x, back.y, 7, 7);


    // рисуем показания
    ctx.strokeStyle = "red";
    let box = {
        leftX: left.x + 3.5 + (data.indication.left * Scale),
        leftY: left.y + 3.5,
        rightX: right.x - 3.5 - (data.indication.right * Scale),
        rightY: right.y + 3.5,
        backX: back.x + 3.5,
        backY: back.y - (data.indication.back * Scale),

        width: data.indication.width_box * Scale,
        length: data.indication.length_box * Scale,
        height: data.indication.height_box * Scale
    };

    if (data.indication.left > 0) {
        ctx.beginPath();
        ctx.moveTo(left.x + 3.5, left.y + 3.5);
        ctx.lineTo(box.leftX, box.leftY);
        ctx.stroke();
    }

    if (data.indication.right > 0) {
        ctx.beginPath();
        ctx.moveTo(right.x + 3.5, right.y + 3.5);
        ctx.lineTo(box.rightX, box.rightY);
        ctx.stroke();
    }

    if (data.indication.back > 0) {
        ctx.beginPath();
        ctx.moveTo(back.x + 3.5, back.y + 3.5);
        ctx.lineTo(box.backX, box.backY);
        ctx.stroke();
    }

    if (box.width > 0 && box.length > 0 && box.height > 0 &&
        data.indication.left > 0 && data.indication.right > 0 && data.indication.back > 0) {
        // рисуем предпологаемый прямоугольник
        ctx.fillStyle = "rgba(166,89,0,0.8)";
        ctx.fillRect(
            box.leftX,
            10,
            box.width,
            box.length);
    }
}

function DrawFront(data) {

    let PlatformHeight = 5;

    let canvas = document.getElementById("canvasFront");
    canvas.width = 648;
    canvas.height = 400;

    let ctx = canvas.getContext("2d");
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // рисуем полоску обозначающая размер СМ
    ctx.font = "8px serif";
    ctx.fillText("1 см", Scale + Scale / 2, 10);

    ctx.beginPath();
    ctx.moveTo(Scale, 15);
    ctx.lineTo(2 * Scale, 15);
    ctx.stroke();

    // рисуем область измерения
    ctx.strokeStyle = "blue";
    ctx.rect(
        canvas.width / 2 - (data.ruler_option.width_max * (Scale / 2)),
        canvas.height - data.ruler_option.top_max * Scale - 5 * Scale,
        data.ruler_option.width_max * Scale,
        data.ruler_option.top_max * Scale,
    );
    ctx.stroke();

    // рисуем платформу весов
    ctx.strokeStyle = "black";
    ctx.fillStyle = "rgba(153,161,163,0.71)";
    ctx.roundRect(
        canvas.width / 2 - (data.scale_platform.width * (Scale / 2)),
        canvas.height - PlatformHeight * Scale,
        data.scale_platform.width * Scale,
        PlatformHeight * Scale,
        5);
    ctx.stroke();
    ctx.fill();

    // рисуем левый датчик
    ctx.fillStyle = "rgb(0,255,0)";
    if (data.indication.left < 0) ctx.fillStyle = "rgb(255,0,0)";
    let left = {
        x: canvas.width / 2 - (data.ruler_option.width_max * (Scale / 2)) - 7,
        y: canvas.height + 7 - PlatformHeight * 2 * Scale
    };
    ctx.fillRect(left.x, left.y, 7, 7);

    // рисуем правый датчик
    ctx.fillStyle = "rgb(0,255,0)";
    if (data.indication.right < 0) ctx.fillStyle = "rgb(255,0,0)";
    let right = {
        x: canvas.width / 2 + (data.ruler_option.width_max * (Scale / 2)),
        y: canvas.height + 10 - PlatformHeight * 2 * Scale
    };
    ctx.fillRect(right.x, right.y, 7, 7);

    // рисуем верхний датчик
    ctx.fillStyle = "rgb(0,255,0)";
    if (data.indication.top < 0) ctx.fillStyle = "rgb(255,0,0)";
    let top = {
        x: canvas.width / 2 - 3.5,
        y: (canvas.height - PlatformHeight * Scale) - (data.ruler_option.top_max * (Scale / 2)) * 2 - 7
    };
    ctx.fillRect(top.x, top.y, 7, 7);

    // рисуем показания
    ctx.strokeStyle = "red";
    let box = {
        leftX: left.x + 3.5 + (data.indication.left * Scale),
        leftY: left.y + 3.5,
        rightX: right.x - 3.5 - (data.indication.right * Scale),
        rightY: right.y + 3.5,
        topX: top.x + 3.5,
        topY: top.y + (data.indication.top * Scale) + 7,

        width: data.indication.width_box * Scale,
        length: data.indication.length_box * Scale,
        height: data.indication.height_box * Scale
    };

    if (data.indication.left > 0) {
        ctx.beginPath();
        ctx.moveTo(left.x + 3.5, left.y + 3.5);
        ctx.lineTo(box.leftX, box.leftY);
        ctx.stroke();
    }

    if (data.indication.right > 0) {
        ctx.beginPath();
        ctx.moveTo(right.x + 3.5, right.y + 3.5);
        ctx.lineTo(box.rightX, box.rightY);
        ctx.stroke();
    }

    if (data.indication.top > 0) {
        ctx.beginPath();
        ctx.moveTo(top.x + 3.5, top.y + 3.5);
        ctx.lineTo(box.topX, box.topY);
        ctx.stroke();
    }

    if (box.width > 0 && box.length > 0 && box.height > 0 &&
        data.indication.left > 0 && data.indication.right > 0 && data.indication.top > 0) {
        // рисуем предпологаемый прямоугольник
        ctx.fillStyle = "rgba(166,89,0,0.8)";
        ctx.fillRect(
            box.leftX,
            canvas.height - PlatformHeight * Scale,
            box.width,
            -box.height);
    }
}

CanvasRenderingContext2D.prototype.roundRect = function (x, y, w, h, r) {
    if (w < 2 * r) r = w / 2;
    if (h < 2 * r) r = h / 2;
    this.beginPath();
    this.moveTo(x + r, y);
    this.arcTo(x + w, y, x + w, y + h, r);
    this.arcTo(x + w, y + h, x, y + h, r);
    this.arcTo(x, y + h, x, y, r);
    this.arcTo(x, y, x + w, y, r);
    this.closePath();
    return this;
};