# mikhmon_v4
# Recode by Irhabi89
# for PHP 8.xx

Udah decode untuk upload ke hosting server. 

bypass lock header localhost.

 Project ini merupakan mikhmon v3 refer to context7, saya ingin anda menganalisi project ini secara menyeluruh pada sisi backendnya pada project ini. dan saya ingi anda membuat planing untuk memigrasikan project ini ke golang dengan menggunakan clean aarchitecture dan sesuai dengan standar project golang. untuk architecture yang perlu anda gunakan pada golang yaitu zap untuk logger,  gorm dengan postgresql, redis untuk caching, go-routeros v3 untuk komunikasi dengan mikoritk, Untuk yang disimpan di postgresql itu hanya konfigurasi routers dan data user untuk login. untuk yang lainnya tetap seperti pada mikhmon. pastikan support multiple mikrotik dan jangan gunakan pooling untuk mendapatkan data secara realtime karena go-routeros ini suuport listen data ke mikrotik menggunakan listenArgs refer to context7.