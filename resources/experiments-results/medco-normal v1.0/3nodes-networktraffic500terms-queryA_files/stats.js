google.maps.__gjsload__('stats', function(_){var P_=function(){this.b=new _.Bl},Q_=function(a){var b=0,c=0;a.b.forEach(function(a){b+=a.Jo;c+=a.ho});return c?b/c:0},R_=function(a,b,c){var d=[];_.gb(a,function(a,c){d.push(c+b+a)});return d.join(c)},S_=function(a){var b={};_.gb(a,function(a,d){b[(0,window.encodeURIComponent)(d)]=(0,window.encodeURIComponent)(a).replace(/%7C/g,"|")});return R_(b,":",",")},T_=function(a,b,c){this.l=b;this.f=a+"/maps/gen_204";this.j=c},U_=function(a,b,c,d){var e={};e.host=window.document.location&&window.document.location.host||
window.location.host;e.v=a;e.r=Math.round(1/b);c&&(e.client=c);d&&(e.key=d);return e},V_=function(a,b,c,d,e){this.m=a;this.C=b;this.l=c;this.f=d;this.j=e;this.b=new _.Bl;this.B=_.Ua()},X_=function(a,b,c,d,e){var f=_.N(_.R,23,500);var g=_.N(_.R,22,2);this.D=a;this.G=b;this.F=f;this.l=g;this.C=c;this.m=d;this.B=e;this.f=new _.Bl;this.b=0;this.j=_.Ua();W_(this)},W_=function(a){window.setTimeout(function(){Y_(a);W_(a)},Math.min(a.F*Math.pow(a.l,a.b),2147483647))},Z_=function(a,b,c){a.f.set(b,c)},Y_=function(a){var b=
U_(a.G,a.C,a.m,a.B);b.t=a.b+"-"+(_.Ua()-a.j);a.f.forEach(function(a,d){a=a();0<a&&(b[d]=Number(a.toFixed(2))+(_.Cm()?"-if":""))});a.D.b({ev:"api_snap"},b);++a.b},$_=function(a,b,c,d,e){this.B=a;this.C=b;this.m=c;this.j=d;this.l=e;this.f={};this.b=[]},a0=function(a,b,c,d){this.j=a;_.G.bind(this.j,"set_at",this,this.l);_.G.bind(this.j,"insert_at",this,this.l);this.B=b;this.C=c;this.m=d;this.f=0;this.b={};this.l()},b0=function(){this.j=_.O(_.R,6);this.C=_.zf();this.b=new T_(_.xg[35]?_.O(_.Af(_.R),11):
_.pw,_.mj,window.document);new a0(_.$i,(0,_.p)(this.b.b,this.b),_.Ff,!!this.j);var a=_.O(new _.tf(_.R.data[3]),1);this.D=a.split(".")[1]||a;this.G={};this.B={};this.m={};this.F={};this.fa=_.N(_.R,0,1);_.lj&&(this.l=new X_(this.b,this.D,this.fa,this.j,this.C))};P_.prototype.f=function(a,b,c){this.b.set(_.Yc(a),{Jo:b,ho:c})};
T_.prototype.b=function(a,b){b=b||{};var c=_.ik().toString(36);b.src="apiv3";b.token=this.l;b.ts=c.substr(c.length-6);a.cad=S_(b);a=R_(a,"=","&");a=this.f+"?target=api&"+a;this.j.createElement("img").src=a;(b=_.pb.__gm_captureCSI)&&b(a)};
V_.prototype.D=function(a,b){b=_.m(b)?b:1;this.b.isEmpty()&&window.setTimeout((0,_.p)(function(){var a=U_(this.C,this.l,this.f,this.j);a.t=_.Ua()-this.B;var b=this.b;_.Cl(b);for(var e={},f=0;f<b.b.length;f++){var g=b.b[f];e[g]=b.H[g]}_.Kz(a,e);this.b.clear();this.m.b({ev:"api_maprft"},a)},this),500);b=this.b.get(a,0)+b;this.b.set(a,b)};$_.prototype.D=function(a){this.f[a]||(this.f[a]=!0,this.b.push(a),2>this.b.length&&_.gA(this,this.G,500))};
$_.prototype.G=function(){for(var a=U_(this.C,this.m,this.j,this.l),b=0,c;c=this.b[b];++b)a[c]="1";a.hybrid=+_.gm();this.b.length=0;this.B.b({ev:"api_mapft"},a)};a0.prototype.l=function(){for(var a;a=this.j.removeAt(0);){var b=a.Kn;a=a.timestamp-this.C;++this.f;this.b[b]||(this.b[b]=0);++this.b[b];if(20<=this.f&&!(this.f%5)){var c={};c.s=b;c.sr=this.b[b];c.tr=this.f;c.te=a;c.hc=this.m?"1":"0";this.B({ev:"api_services"},c)}}};b0.prototype.T=function(a){a=_.Yc(a);this.G[a]||(this.G[a]=new $_(this.b,this.D,this.fa,this.j,this.C));return this.G[a]};b0.prototype.S=function(a){a=_.Yc(a);this.B[a]||(this.B[a]=new V_(this.b,this.D,_.N(_.R,0,1),this.j,this.C));return this.B[a]};b0.prototype.f=function(a){if(this.l){this.m[a]||(this.m[a]=new _.ZA,Z_(this.l,a,function(){return b.Ya()}));var b=this.m[a];return b}};
b0.prototype.N=function(a){if(this.l){this.F[a]||(this.F[a]=new P_,Z_(this.l,a,function(){return Q_(b)}));var b=this.F[a];return b}};var c0=new b0;_.ke("stats",c0);});
