import time

from selenium import webdriver
from telegram.ext import Updater

updater = Updater('TELEGRAM_BOT_TOKEN')
def get_data():
    options = webdriver.ChromeOptions()
    options.add_argument("--headless")
    browser = webdriver.Chrome("chromedriver.exe",desired_capabilities=options.to_capabilities())
    try:
        browser.get("https://scolkg.com/")
        a = browser.find_element_by_xpath('//*[@id="app_coinboard"]/div[2]/table/tbody/tr[1]/td[6]').text
        browser.quit()
        return float(a[a.find('(') + 1:len(a) - 2])
    except:
        updater.bot.sendMessage(TELEGRAM_CHAT_CODE, "김프 사이트에서 데이터를 읽어오고 편집하는데 실패하였습니다.")
while True:
    a = get_data()
    if a >= 15:
        for d in range(0, 15):
            updater.bot.sendMessage(TELEGRAM_CHAT_CODE, "김치 프리미엄이 15%를 상회하였습니다.")
    print('ok')
    time.sleep(300)
