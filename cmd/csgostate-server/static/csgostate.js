function poll() {
    setTimeout(function () {
        GetData();
    }, 1000);
}

function GetData() {
    console.log("Pollign");
    jQuery.ajax({
        url: "api/players",
        type: "GET",
        dataType: "json",
        contentType: 'application/json; charset=utf-8',
        success: function(resultData) {
            console.log("success");
            let andy = resultData["76561197993200126"];
            if (!andy) {
                console.log("no andy", resultData);
                return
            }
            let LatestState = andy.LatestState;
            if (!LatestState) {
                console.log("no state", andy);
                return
            }
            let player = LatestState.player;
            if (!player) {
                console.log("no player", LatestState);
                return
            }
             $("#players").html(JSON.stringify(player));
        },
        error : function(xhr, textStatus, errorThrown) {
            console.log("error", textStatus);
        },
        complete: function() {
            poll();
        },
        timeout: 0,
    });
}

poll();

/*
{"76561197993200126":{"LatestState":{"provider":{"steamid":"76561197993200126","timestamp":1612128081},"player":{"activity":"playing","state":{"health":100,"armor":100,"helmet":true,"flashed":0,"smoked":0,"burning":0,"money":0,"round_kills":0,"round_killhs":0,"equip_value":5950},"weapons":{"weapon_0":{"name":"weapon_knife_t","paintkit":"default","type":"Knife","ammo_clip":0,"ammo_clip_max":0,"ammo_reserve":0,"state":"holstered"},"weapon_1":{"name":"weapon_glock","paintkit":"aq_glock18_flames_blue","type":"Pistol","ammo_clip":20,"ammo_clip_max":20,"ammo_reserve":120,"state":"holstered"},"weapon_2":{"name":"weapon_awp","paintkit":"gs_awp_phobos","type":"SniperRifle","ammo_clip":10,"ammo_clip_max":10,"ammo_reserve":30,"state":"active"}}},"round":{"phase":"freezetime","bomb":""},"RawJson":"{\n\t\"provider\": {\n\t\t\"name\": \"Counter-Strike: Global Offensive\",\n\t\t\"appid\": 730,\n\t\t\"version\": 13779,\n\t\t\"steamid\": \"76561197993200126\",\n\t\t\"timestamp\": 1612128081\n\t},\n\t\"player\": {\n\t\t\"steamid\": \"76561197993200126\",\n\t\t\"clan\": \"VUKUKARSK\",\n\t\t\"name\": \"Andy\
 */
function printPlayerHP() {

}