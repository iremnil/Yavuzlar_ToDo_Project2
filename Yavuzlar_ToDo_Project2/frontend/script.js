"use strict";

const inputGiris = document.getElementById("gorevInput");
const butonEkle = document.getElementById("ekleBtn");
const listeUl = document.getElementById("gorevListesi");
const inputArama = document.getElementById("aramaInput");
const editModal = document.getElementById("editModal");
const modalInput = document.getElementById("modalInput");
const modalKaydetBtn = document.getElementById("modalKaydetBtn");

const deleteModal = document.getElementById("deleteModal");
const modalSilBtn = document.getElementById("modalSilBtn");

let islemYapilacakId = null;
// Sayfa yüklendiğinde görevleri çek
document.addEventListener("DOMContentLoaded", gorevleriGetir);
async function gorevleriGetir() {
    try {
        const response = await fetch('/api/v1/todos');
        
        if (response.status === 401) {
            window.location.replace("/login.html");
            return;
        }
        if (response.ok) {
            const data = await response.json();
            listeyiYenile(data); 
            //dınamık username
            const userResponse = await fetch('/api/v1/me');
            if (userResponse.ok) {
                const userData = await userResponse.json();
                document.getElementById("kullaniciAdi").innerText = userData.username;
            }
            
        }
    } catch (error) {
        console.error("Görevler çekilemedi veya sunucuya ulaşılamadı:", error);
        // Sunucuya hiç ulaşılamazsa da guvenlık ıcın logıne at
        window.location.replace("/login.html"); 
    }
}
// Add buton
butonEkle.addEventListener("click", async () => {
    const metin = inputGiris.value.trim();
    if (metin !== "") {
        const response = await fetch('/api/v1/todos', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title: metin, status: "Yapılacak" }) 
        });
        
        if (response.ok) {
            inputGiris.value = "";
            gorevleriGetir(); 
        }
    }
});

//arama cubugu
inputArama.addEventListener("input", function() {
    const arananKelime = this.value.toLowerCase(); 
    const listeElemanlari = document.querySelectorAll("#gorevListesi li"); 

    listeElemanlari.forEach(li => {
        const gorevMetni = li.querySelector(".gorev-metni").textContent.toLowerCase();
        if (gorevMetni.includes(arananKelime)) {
            li.style.display = ""; 
        } else {
            li.style.display = "none"; 
        }
    });
});
// Liste yenileme
function listeyiYenile(data) {
    listeUl.innerHTML = "";
    data.forEach((is, index) => {
        const li = document.createElement("li");
        const durumMetni = is.status || "Yapılacak"; 
        let durumSinifi = "durum-yapilacak";
        if (durumMetni === "Devam Ediyor") durumSinifi = "durum-devam";
        if (durumMetni === "Tamamlandı") durumSinifi = "durum-tamamlandi";
        li.innerHTML = `
            <span class="sira-no">${index + 1}</span>
            <span class="gorev-metni">${is.title}</span> 
            <button class="status-btn ${durumSinifi}" onclick="durumDegistir('${is.id}', '${is.title.replace(/'/g, "\\'")}', '${durumMetni}')">${durumMetni}</button>
            <div class="btn-grup">
                <button class="edit-btn" onclick="guncelle('${is.id}', '${is.title.replace(/'/g, "\\'")}')">Düzenle</button>
                <button class="delete-btn" onclick="sil('${is.id}')">Sil</button>
            </div>
        `;
        listeUl.appendChild(li);
    });
    ilerlemeGuncelle(data);
}

function guncelle(id, eskiTitle) {   //edit
    islemYapilacakId = id;                  
    modalInput.value = eskiTitle;           
    editModal.style.display = "flex";       
}

function modalKapat() {
    islemYapilacakId = null;
    modalInput.value = "";
    editModal.style.display = "none";       
}
modalKaydetBtn.addEventListener("click", async () => {
    const yeniBaslik = modalInput.value.trim();
    if (islemYapilacakId && yeniBaslik !== "") {
        const response = await fetch(`/api/v1/todos?id=${islemYapilacakId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title: yeniBaslik }) 
        });
        
        if (response.ok) {
            modalKapat();       
            gorevleriGetir();   
        }
    }
});
function sil(id) {
    islemYapilacakId = id;                  
    deleteModal.style.display = "flex";     
}
function deleteModalKapat() {
    islemYapilacakId = null;
    deleteModal.style.display = "none";     
}

modalSilBtn.addEventListener("click", async () => {
    if (islemYapilacakId) {
        const response = await fetch(`/api/v1/todos?id=${islemYapilacakId}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            deleteModalKapat(); 
            gorevleriGetir();   
        }
    }
});

async function performLogout() {  //Cıkıs
    try {
        const response = await fetch('/api/v1/logout', { method: 'POST' });
        
        if (response.ok) {
            window.location.replace("/login.html"); 
        }
    } catch (error) {
        console.error("Çıkış hatası:", error);
    }
}

//durum degıstırme
async function durumDegistir(id, mevcutTitle, mevcutDurum) {
    let yeniDurum = "Yapılacak"; 

    if (mevcutDurum === "Yapılacak") {
        yeniDurum = "Devam Ediyor";
    } else if (mevcutDurum === "Devam Ediyor") {
        yeniDurum = "Tamamlandı";
    } else if (mevcutDurum === "Tamamlandı") {
        yeniDurum = "Yapılacak"; 
    } else {
        yeniDurum = "Devam Ediyor";
    }
    const response = await fetch(`/api/v1/todos?id=${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: mevcutTitle, status: yeniDurum }) 
    });

    if (response.ok) {
        gorevleriGetir(); 
    }
}
//ZAMAN VE KARŞILAMA WIDGET
function saatiGuncelle() {
    const simdi = new Date();
    // Tarih Ayar
    const tarihSecenekleri = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
    const formatliTarih = simdi.toLocaleDateString('tr-TR', tarihSecenekleri);
    // Saat Ayar
    const formatliSaat = simdi.toLocaleTimeString('tr-TR', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
    document.getElementById("tarih").innerText = formatliTarih;
    document.getElementById("saat").innerText = formatliSaat;
}
// Saati her 1 saniyede güncelle
setInterval(saatiGuncelle, 1000);
saatiGuncelle(); // Sayfa açılır açılmaz
//İLERLEME ÇEMBERİNİ GÜNCELLEME
let mevcutYuzde = 0;
let animasyonTimer = null;
function ilerlemeGuncelle(data) {
    const toplamGorev = data.length;
    const daire = document.getElementById("ilerlemeDairesi");            //animasyon
    const yuzdeMetni = document.getElementById("ilerlemeYuzde");

    let hedefYuzde = 0;
    if (toplamGorev > 0) {
        const tamamlananGorev = data.filter(gorev => gorev.status === "Tamamlandı").length;
        hedefYuzde = Math.round((tamamlananGorev / toplamGorev) * 100);
    }
    // Eğer zaten hedefteyse boşuna çalışma
    if (mevcutYuzde === hedefYuzde) return;
    // Önceki yarım kalan animasyonu durdur (hızlı hızlı tıklanırsa çakışmasın)
    if (animasyonTimer) clearInterval(animasyonTimer);
    // İlerleme yönüne göre animasyon renk
    let animasyonRengi = "#3498db"; // Standart mavi
    if (hedefYuzde > mevcutYuzde) {
        animasyonRengi = "#55efc4"; // Artıyorsa yeşil 
    } else if (hedefYuzde < mevcutYuzde) {
        animasyonRengi = "#ff7675"; // Azalıyorsa kırmızı
    }
    //her 15 ms yüzdeyi 1 birim kaydır
    animasyonTimer = setInterval(() => {
        // Hedefe gelındıyse durdur
        if (mevcutYuzde === hedefYuzde) {
            clearInterval(animasyonTimer);
            // Animasyon bitince, eğer %100 değilse rengi tekrar sakin maviye döndür
            if (mevcutYuzde !== 100 && mevcutYuzde !== 0) {
                const derece = (mevcutYuzde / 100) * 360;
                daire.style.background = `conic-gradient(#3498db ${derece}deg, #eee 0deg)`;
            } else if (mevcutYuzde === 100) {
                 daire.style.background = `conic-gradient(#55efc4 360deg, #eee 0deg)`; // 100'de yeşil kalsın
            } else if (mevcutYuzde === 0) {
                 daire.style.background = `conic-gradient(#eee 360deg, #eee 0deg)`; // 0'da gri kalsın
            }
            return;
        }
        // Yüzdeyi 1'er 1'er artır veya azalt
        if (mevcutYuzde < hedefYuzde) {
            mevcutYuzde++;
        } else {
            mevcutYuzde--;
        }
        // Çemberi yeni yüzdeye ve renge göre tekrar çiz
        const derece = (mevcutYuzde / 100) * 360;
        daire.style.background = `conic-gradient(${animasyonRengi} ${derece}deg, #eee 0deg)`;
        yuzdeMetni.innerText = `%${mevcutYuzde}`;

    }, 15); // Hız ayarı 
}