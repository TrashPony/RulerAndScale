// библиотека для работы с I²C хабом
#include <TroykaI2CHub.h>
// библиотека для управления лазером
#include "Adafruit_VL53L0X.h"

// объект для работы с хабом адрес по умолчанию 0x70
TroykaI2CHub splitter;

// адрес устройства лазера один для всех подсетей
#define LOX2_ADDRESS 0x29

// Кнопки
#define BUTTON 8

// Дальномер
#define LEFT_PING_LAN  0
#define RIGHT_PING_LAN 1
#define TOP_PING_LAN   2
#define BACK_PING_LAN  3

// Диоды
#define RED_LED_PIN  9
#define GREEN_LED_PIN  5

#define TOP_MAX          77
#define WIDTH_MAX        87
#define LENGTH_MAX       46

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

boolean debug = true;

Adafruit_VL53L0X lox = Adafruit_VL53L0X();

void setup()
{
  // подключаем лазер
  lox.begin();

  pinMode(RED_LED_PIN, OUTPUT);
  pinMode(GREEN_LED_PIN, OUTPUT);
  pinMode(BUTTON, INPUT);

  Serial.begin(115200);
  // ждем пока откроется сериал порт
  while(!Serial) {}

  // начало работы с I²C хабом
  splitter.begin();

  // ждём одну секунду
  delay(1000);
}

void loop()
{
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

int getDistance(int pin) {

  // pin - указываем номер сети для лазера откуда брать данные
  splitter.setBusChannel(pin);
  
  return getIndication();
}

int getIndication() {
  
  VL53L0X_RangingMeasurementData_t measure;
  lox.rangingTest(&measure, false);

  // проверка на доступность устройства 
  Wire.beginTransmission(LOX2_ADDRESS);
  byte state = Wire.endTransmission();

  float dist = round((measure.RangeMilliMeter)/10);
  int distInt = int(dist);
  
  if (measure.RangeStatus != 4 && state == 0) {
    return distInt;
  } else {
    return -1;
  }   
}

void Indication() {

  int right = getDistance(RIGHT_PING_LAN);
  int left = getDistance(LEFT_PING_LAN);
  int top =  getDistance(TOP_PING_LAN);
  int back = getDistance(BACK_PING_LAN);

  if (right > 0 && left > 0 && top > 0 && back > 0) {

    widthBox = WIDTH_MAX - (right + left);
    heightBox = TOP_MAX - top;
    lengthBox = LENGTH_MAX - back;

    passWidthBox = PassedIndication(lastWidthBox, widthBox, queueIndicationsWidth);
    passHeightBox = PassedIndication(lastHeightBox, heightBox, queueIndicationsHeight);
    passLengthBox = PassedIndication(lastLengthBox, lengthBox, queueIndicationsLength);

  } else {
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
