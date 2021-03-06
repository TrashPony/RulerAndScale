#include <NewPing.h>

// Кнопки
#define BUTTON 8

// Дальномер
#define LEFT_PING_PIN  13
#define TOP_PING_PIN   12
#define BACK_PING_PIN  11
#define RIGHT_PING_PIN 10

// Диоды
#define RED_LED_PIN  9
#define GREEN_LED_PIN  5

#define TOP_MAX          87
#define WIDTH_MAX        103
#define LENGTH_MAX       61

boolean onlyWeight = false;

boolean passWidthBox = false;
boolean passHeightBox = false;
boolean passLengthBox = false;

int lastWidthBox;
int lastHeightBox;
int lastLengthBox;

int widthBox;
int heightBox;
int lengthBox;

int queueIndicationsWidth = 0;
int queueIndicationsHeight = 0;
int queueIndicationsLength = 0;

boolean debug = false;

NewPing rightSonar(RIGHT_PING_PIN, RIGHT_PING_PIN, WIDTH_MAX);
NewPing leftSonar(LEFT_PING_PIN, LEFT_PING_PIN, WIDTH_MAX);
NewPing topSonar(TOP_PING_PIN, TOP_PING_PIN, TOP_MAX);
NewPing backSonar(BACK_PING_PIN, BACK_PING_PIN, LENGTH_MAX);


void setup() {
  Serial.begin(115200);
  pinMode(RED_LED_PIN, OUTPUT);
  pinMode(GREEN_LED_PIN, OUTPUT);
  pinMode(BUTTON, INPUT);
}

void loop() {

  if (digitalRead(BUTTON) == LOW) {
    if (onlyWeight) {
      onlyWeight = false;
      digitalWrite(RED_LED_PIN, HIGH);
      digitalWrite(GREEN_LED_PIN, LOW);
    } else {
      onlyWeight = true;
      digitalWrite(RED_LED_PIN, LOW);
      digitalWrite(GREEN_LED_PIN, HIGH);
    }
    delay(500);
  }

  Indication();

  if (Serial.available()) {

    byte incomingByte = Serial.read();

    if (incomingByte == 0x95) {
      Serial.write(0x7F);
      Serial.write(0x7F);
      Serial.write(0x7F);
      Serial.write(0x7F);
      Serial.write(0x7F);
      Serial.flush();
    }

    // запрос ширины
    if (incomingByte == 0x88) {
      Serial.write(0x2D);

      if (onlyWeight) {
        Serial.write(0x7A);
      }
      else {
        if (passWidthBox) {
          Serial.write(0x7F);
        }
        else {
          Serial.write(0x7E);
        }
      }

      Serial.write(0x0B);
      Serial.write(widthBox);
      Serial.write(0x7B);
    }

    // запрос высоты
    if (incomingByte == 0x99) {
      Serial.write(0x2D);

      if (onlyWeight) {
        Serial.write(0x7A);
      }
      else {
        if (passHeightBox) {
          Serial.write(0x7F);
        }
        else {
          Serial.write(0x7E);
        }
      }

      Serial.write(0x16);
      Serial.write(heightBox);
      Serial.write(0x7B);
    }

    // запрос длины
    if (incomingByte == 0x77) {
      Serial.write(0x2D);

      if (onlyWeight) {
        Serial.write(0x7A);
      }
      else {
        if (passLengthBox) {
          Serial.write(0x7F);
        }
        else {
          Serial.write(0x7E);
        }
      }

      Serial.write(0x21);
      Serial.write(lengthBox);
      Serial.write(0x7B);
    }

    if (incomingByte == 0x66) {
      if (onlyWeight) {
        digitalWrite(GREEN_LED_PIN, HIGH);
        digitalWrite(RED_LED_PIN, LOW);
      }
      else {
        digitalWrite(RED_LED_PIN, HIGH);
        digitalWrite(GREEN_LED_PIN, LOW);
      }
    }

    if (incomingByte == 0x55) {
      digitalWrite(RED_LED_PIN, LOW);
      digitalWrite(GREEN_LED_PIN, LOW);
    }
  }
}

int SearchAvg (int indications[], int countIndications) {

  int result = 0;
  int maxCount = 0;

  for (int i = 0; i < countIndications; i++) {
    int count = 0;
    for (int j = 0; j < countIndications; j++) {
      if (indications[i] == indications[j]) {
        count++;
      }

      if (maxCount < count) {
        maxCount = count;
        result = i;
      }
    }
  }

  return indications[result];
}

void Indication() {
  int countIndications = 5;

  int rightIndications[countIndications];
  int leftIndications[countIndications];
  int topIndications[countIndications];
  int backIndications[countIndications];


  for (int i = 0; i < countIndications; i++) {

    rightIndications[i] = rightSonar.ping_cm();
    leftIndications[i] = leftSonar.ping_cm();
    topIndications[i] = topSonar.ping_cm();
    backIndications[i] = backSonar.ping_cm();

  }

  int right = SearchAvg(rightIndications, countIndications);
  int left = SearchAvg(leftIndications, countIndications);
  int top = SearchAvg(topIndications, countIndications);
  int back = SearchAvg(backIndications, countIndications);

  if (right > 0 && left > 0 && top > 0 && back > 0) {

    widthBox = WIDTH_MAX - (right + left);
    heightBox = TOP_MAX - top;
    lengthBox = LENGTH_MAX - back;

    passWidthBox = PassedIndication(lastWidthBox, widthBox, queueIndicationsWidth);
    passHeightBox = PassedIndication(lastHeightBox, heightBox, queueIndicationsHeight);
    passLengthBox = PassedIndication(lastLengthBox, lengthBox, queueIndicationsLength);

  }
  else {
    passWidthBox = false;
    passHeightBox = false;
    passLengthBox = false;

    lastWidthBox = 0;
    lastHeightBox = 0;
    lastLengthBox = 0;

    widthBox = 0;
    heightBox = 0;
    lengthBox = 0;
  }

  if (debug) {
    Serial.print("Right_ping: ");
    Serial.print(right);
    Serial.println("cm");

    Serial.print("Left_ping: ");
    Serial.print(left);
    Serial.println("cm");

    Serial.print("Top_ping: ");
    Serial.print(top);
    Serial.println("cm");

    Serial.print("Back_ping: ");
    Serial.print(back);
    Serial.println("cm");

    Serial.println("///////////////////////////////////////");

    Serial.print("widthBox ");
    Serial.print(passWidthBox);
    Serial.print(" ");
    Serial.println(widthBox);

    Serial.print("heightBox ");
    Serial.print(passHeightBox);
    Serial.print(" ");
    Serial.println(heightBox);

    Serial.print("lengthBox ");
    Serial.print(passLengthBox);
    Serial.print(" ");
    Serial.println(lengthBox);

    Serial.println("///////////////////////////////////////");
    delay(2500);
  }
}

boolean PassedIndication (int &lastIndication, int &indication, int &queueIndications) {

  if (debug) {
    Serial.print("queueIndications ");
    Serial.println(queueIndications);

    Serial.print("lastIndication ");
    Serial.println(lastIndication);

    Serial.print("indication ");
    Serial.println(indication);
  }

  if (indication > 0) {
    if ((lastIndication - 1) <= indication && indication <= (lastIndication + 1)) {
      if (queueIndications >= 5) {
        return true;
      }
      else {
        queueIndications++;
        return false;
      }
    }
    else {
      queueIndications = 0;
      lastIndication = indication;
      return false;
    }
  }
  else {
    queueIndications = 0;
    return false;
  }
}
