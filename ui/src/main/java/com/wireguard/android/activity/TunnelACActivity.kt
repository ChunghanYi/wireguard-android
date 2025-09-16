/*
 * Copyright © 2025 Slowboot(chunghan.yi@gmail.com). All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package com.wireguard.android.activity

import android.os.Bundle
import com.wireguard.android.R
import com.wireguard.android.model.ObservableTunnel

/**
 * Standalone activity for auto connect.
 */
class TunnelACActivity : BaseActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.tunnel_ac_activity)
    }

    override fun onSelectedTunnelChanged(oldTunnel: ObservableTunnel?, newTunnel: ObservableTunnel?): Boolean {
        finish()
        return true
    }
}
