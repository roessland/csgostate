package main

const OpenIDProvider = "https://steamcommunity.com/openid"
// https://steamcommunity.com/openid/login

// login -> show "login through steam button"
// <form action="https://steamcommunity.com/openid/login" method="post">
//<input type="hidden" name="openid.identity" value="http://specs.openid.net/auth/2.0/identifier_select">
//<input type="hidden" name="openid.claimed_id" value="http://specs.openid.net/auth/2.0/identifier_select">
//<input type="hidden" name="openid.ns" value="http://specs.openid.net/auth/2.0">
//<input type="hidden" name="openid.mode" value="checkid_setup">
//<input type="hidden" name="openid.realm" value="https://steamdb.info/">
//<input type="hidden" name="openid.return_to" value="https://steamdb.info/login/">
//<button type="submit" class="btn btn-steam" id="js-sign-in">Sign in through Steam</button>
//</form>
// https://csgostate.roessland.com/login?
//	openid.ns=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0
//	&openid.mode=id_res
//	&openid.op_endpoint=https%3A%2F%2Fsteamcommunity.com%2Fopenid%2Flogin
//	&openid.claimed_id=https%3A%2F%2Fsteamcommunity.com%2Fopenid%2Fid%2F76561197993200126
//	&openid.identity=https%3A%2F%2Fsteamcommunity.com%2Fopenid%2Fid%2F76561197993200126
//	&openid.return_to=https%3A%2F%2Fcsgostate.roessland.com%2Flogin
//	&openid.response_nonce=2021-02-01T22%3A11%3A05Z0nWH2sZk%2B2S%2Fqy%2FvhtCW87Fsekc%3D
//	&openid.assoc_handle=1234567890
//	&openid.signed=signed%2Cop_endpoint%2Cclaimed_id%2Cidentity%2Creturn_to%2Cresponse_nonce%2Cassoc_handle
//	&openid.sig=%2Fn4qdELlWvfmq2fSCL9%2FOdJfvjw%3D

/*
To verify the user, make a call from your backend to https://steamcommunity.com/openid/login copying every query string parameter from that response with one exception: replace &openid.mode=id_res with &openid.mode=check_authentication. So the final call will be to this URL:
 */