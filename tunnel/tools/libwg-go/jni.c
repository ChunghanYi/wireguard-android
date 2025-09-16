/* SPDX-License-Identifier: Apache-2.0
 *
 * Copyright © 2017-2021 Jason A. Donenfeld <Jason@zx2c4.com>. All Rights Reserved.
 */

#include <jni.h>
#include <stdlib.h>
#include <string.h>

struct go_string { const char *str; long n; };
extern int wgTurnOn(struct go_string ifname, int tun_fd, struct go_string settings);
extern void wgTurnOff(int handle);
extern int wgGetSocketV4(int handle);
extern int wgGetSocketV6(int handle);
extern char *wgGetConfig(int handle);
extern char *wgVersion();
// AutoConnect --
extern char *acTurnOn(struct go_string serverip, struct go_string port,
		struct go_string privatekey, struct go_string publickey);
extern int acTurnOff(struct go_string serverip, struct go_string port,
		struct go_string publickey);
// -- -- --

JNIEXPORT jint JNICALL Java_com_wireguard_android_backend_GoBackend_wgTurnOn(JNIEnv *env, jclass c, jstring ifname, jint tun_fd, jstring settings)
{
	const char *ifname_str = (*env)->GetStringUTFChars(env, ifname, 0);
	size_t ifname_len = (*env)->GetStringUTFLength(env, ifname);
	const char *settings_str = (*env)->GetStringUTFChars(env, settings, 0);
	size_t settings_len = (*env)->GetStringUTFLength(env, settings);
	int ret = wgTurnOn((struct go_string){
		.str = ifname_str,
		.n = ifname_len
	}, tun_fd, (struct go_string){
		.str = settings_str,
		.n = settings_len
	});
	(*env)->ReleaseStringUTFChars(env, ifname, ifname_str);
	(*env)->ReleaseStringUTFChars(env, settings, settings_str);
	return ret;
}

JNIEXPORT void JNICALL Java_com_wireguard_android_backend_GoBackend_wgTurnOff(JNIEnv *env, jclass c, jint handle)
{
	wgTurnOff(handle);
}

JNIEXPORT jint JNICALL Java_com_wireguard_android_backend_GoBackend_wgGetSocketV4(JNIEnv *env, jclass c, jint handle)
{
	return wgGetSocketV4(handle);
}

JNIEXPORT jint JNICALL Java_com_wireguard_android_backend_GoBackend_wgGetSocketV6(JNIEnv *env, jclass c, jint handle)
{
	return wgGetSocketV6(handle);
}

JNIEXPORT jstring JNICALL Java_com_wireguard_android_backend_GoBackend_wgGetConfig(JNIEnv *env, jclass c, jint handle)
{
	jstring ret;
	char *config = wgGetConfig(handle);
	if (!config)
		return NULL;
	ret = (*env)->NewStringUTF(env, config);
	free(config);
	return ret;
}

JNIEXPORT jstring JNICALL Java_com_wireguard_android_backend_GoBackend_wgVersion(JNIEnv *env, jclass c)
{
	jstring ret;
	char *version = wgVersion();
	if (!version)
		return NULL;
	ret = (*env)->NewStringUTF(env, version);
	free(version);
	return ret;
}

// AutoConnect --
JNIEXPORT jstring JNICALL Java_com_wireguard_android_backend_GoBackend_acTurnOn(JNIEnv *env, jclass c,
		jstring serverip, jstring port, jstring privatekey, jstring publickey)
{
	const char *serverip_str = (*env)->GetStringUTFChars(env, serverip, 0);
	size_t serverip_len = (*env)->GetStringUTFLength(env, serverip);
	const char *port_str = (*env)->GetStringUTFChars(env, port, 0);
	size_t port_len = (*env)->GetStringUTFLength(env, port);
	const char *privatekey_str = (*env)->GetStringUTFChars(env, privatekey, 0);
	size_t privatekey_len = (*env)->GetStringUTFLength(env, privatekey);
	const char *publickey_str = (*env)->GetStringUTFChars(env, publickey, 0);
	size_t publickey_len = (*env)->GetStringUTFLength(env, publickey);

	char *config = acTurnOn(
		(struct go_string){
			.str = serverip_str,
			.n = serverip_len
		},
		(struct go_string){
			.str = port_str,
			.n = port_len
		},
		(struct go_string){
			.str = privatekey_str,
			.n = privatekey_len
		},
		(struct go_string){
			.str = publickey_str,
			.n = publickey_len
		});

	if (!config) {
		(*env)->ReleaseStringUTFChars(env, serverip, serverip_str);
		(*env)->ReleaseStringUTFChars(env, port, port_str);
		(*env)->ReleaseStringUTFChars(env, privatekey, privatekey_str);
		(*env)->ReleaseStringUTFChars(env, publickey, publickey_str);
		return NULL;
	} else {
		jstring ret = (*env)->NewStringUTF(env, config);
		free(config);

		(*env)->ReleaseStringUTFChars(env, serverip, serverip_str);
		(*env)->ReleaseStringUTFChars(env, port, port_str);
		(*env)->ReleaseStringUTFChars(env, privatekey, privatekey_str);
		(*env)->ReleaseStringUTFChars(env, publickey, publickey_str);
		return ret;
	}
}

JNIEXPORT jint JNICALL Java_com_wireguard_android_backend_GoBackend_acTurnOff(JNIEnv *env, jclass c,
		jstring serverip, jstring port, jstring publickey)
{
	const char *serverip_str = (*env)->GetStringUTFChars(env, serverip, 0);
	size_t serverip_len = (*env)->GetStringUTFLength(env, serverip);
	const char *port_str = (*env)->GetStringUTFChars(env, port, 0);
	size_t port_len = (*env)->GetStringUTFLength(env, port);
	const char *publickey_str = (*env)->GetStringUTFChars(env, publickey, 0);
	size_t publickey_len = (*env)->GetStringUTFLength(env, publickey);

	int ret = acTurnOff(
		(struct go_string){
			.str = serverip_str,
			.n = serverip_len
		},
		(struct go_string){
			.str = port_str,
			.n = port_len
		},
		(struct go_string){
			.str = publickey_str,
			.n = publickey_len
		});

	(*env)->ReleaseStringUTFChars(env, serverip, serverip_str);
	(*env)->ReleaseStringUTFChars(env, port, port_str);
	(*env)->ReleaseStringUTFChars(env, publickey, publickey_str);

	return ret;
}
// -- -- --
