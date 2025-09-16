/*
 * Copyright © 2025 Slowboot(chunghan.yi@gmail.com). All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package com.wireguard.android.fragment

import android.os.Bundle
import android.util.Log
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.lifecycle.lifecycleScope
import com.wireguard.android.backend.Tunnel
import com.wireguard.android.Application
import com.wireguard.android.databinding.TunnelAcFragmentBinding
import com.wireguard.android.Application.Companion.get
import com.wireguard.android.Application.Companion.getBackend
import com.wireguard.config.Config
import com.wireguard.android.configStore.ConfigStore
import com.wireguard.android.model.ObservableTunnel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext

/**
 * A simple [Fragment] subclass.
 * Use the [TunnelACFragment.newInstance] factory method to
 * create an instance of this fragment.
 */
class TunnelACFragment : Fragment() {
    private var binding: TunnelAcFragmentBinding? = null
    val manager = Application.getTunnelManager()
    var actunnel: ObservableTunnel? = null	//from TunnelEditorFragment.kt

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
    }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        super.onCreateView(inflater, container, savedInstanceState)
        binding = TunnelAcFragmentBinding.inflate(inflater, container, false)
        binding?.apply {
            // Access your button directly through the binding object
            acStartButton.setOnClickListener {
                lifecycleScope.launch {
                    acTunnelConfig(acServerIp.text.toString(), acServerPort.text.toString())
                }
            }
        }
        return binding?.root
    }

    suspend fun acTunnelConfig(serverIp: String, serverPort: String) {
        var autoConfig: Config? = null
        var configStore: ConfigStore

        if (!manager.getTunnels().containsKey("wg0")) {
            autoConfig = getBackend().setAC(serverIp, serverPort, Tunnel.State.UP)
            if (autoConfig != null) {
                if (actunnel != null) {
                    manager.delete(actunnel!!)
                }
                actunnel = manager.create("wg0", autoConfig)
                actunnel!!.setConfigAsync(autoConfig!!)
                actunnel!!.setNameAsync("wg0")

                val ctx = activity ?: Application.get()
                val message = "OK, Auto Connection is started."
                Toast.makeText(ctx, message, Toast.LENGTH_SHORT).show()

                parentFragmentManager.beginTransaction()
                    .remove(this) // 'this' refers to the current Fragment instance
                    .commit()
            } else {
                val ctx = activity ?: Application.get()
                val message = "Oops, Auto Connection is failed."
                Toast.makeText(ctx, message, Toast.LENGTH_SHORT).show()
            }
        } else {
            val ctx = activity ?: Application.get()
            val message = "You should try this after deleting the wg0 tunnel you set up earlier."
            Toast.makeText(ctx, message, Toast.LENGTH_SHORT).show()
        }
    }

    companion object {
    }
}
