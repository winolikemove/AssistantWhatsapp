package ai

// SalesToolDefinitions contains OpenAI-compatible function-calling tool schemas for sales.
var SalesToolDefinitions = []map[string]interface{}{
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "create_sales_transaction",
			"description": "Catat transaksi penjualan baru. Gunakan ini ketika user ingin menjual barang ke customer.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"item_nama": map[string]interface{}{
						"type":        "string",
						"description": "Nama barang yang dijual",
					},
					"qty": map[string]interface{}{
						"type":        "integer",
						"description": "Jumlah/quantity barang",
					},
					"customer_nama": map[string]interface{}{
						"type":        "string",
						"description": "Nama customer pembeli",
					},
					"catatan": map[string]interface{}{
						"type":        "string",
						"description": "Catatan tambahan (opsional)",
					},
				},
				"required": []string{"item_nama", "qty", "customer_nama"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "add_sales_item",
			"description": "Tambah item baru ke database dengan harga beli tetap dari supplier.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"nama": map[string]interface{}{
						"type":        "string",
						"description": "Nama item/barang",
					},
					"harga_beli": map[string]interface{}{
						"type":        "number",
						"description": "Harga beli dari supplier dalam Rupiah",
					},
					"satuan": map[string]interface{}{
						"type":        "string",
						"description": "Satuan: kg, liter, pcs, dll",
					},
				},
				"required": []string{"nama", "harga_beli"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "add_sales_customer",
			"description": "Tambah customer baru ke database.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"nama": map[string]interface{}{
						"type":        "string",
						"description": "Nama customer",
					},
					"alamat": map[string]interface{}{
						"type":        "string",
						"description": "Alamat lengkap",
					},
					"telepon": map[string]interface{}{
						"type":        "string",
						"description": "Nomor telepon/WhatsApp",
					},
					"jatuh_tempo": map[string]interface{}{
						"type":        "integer",
						"description": "Jatuh tempo dalam hari: 7, 14, atau 30",
					},
					"payment": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"cash", "credit"},
						"description": "Metode pembayaran",
					},
				},
				"required": []string{"nama", "alamat", "jatuh_tempo", "payment"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "set_customer_pricing",
			"description": "Set harga jual khusus untuk item tertentu ke customer tertentu.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"customer_nama": map[string]interface{}{
						"type":        "string",
						"description": "Nama customer",
					},
					"item_nama": map[string]interface{}{
						"type":        "string",
						"description": "Nama item",
					},
					"harga_jual": map[string]interface{}{
						"type":        "number",
						"description": "Harga jual per unit dalam Rupiah",
					},
				},
				"required": []string{"customer_nama", "item_nama", "harga_jual"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "get_profit_report",
			"description": "Mendapatkan laporan keuntungan untuk periode tertentu.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"period": map[string]interface{}{
						"type":        "string",
						"description": "Periode: hari ini, minggu ini, bulan ini",
					},
				},
				"required": []string{"period"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "get_receivable_summary",
			"description": "Mendapatkan ringkasan piutang dari customer.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties":  map[string]interface{}{},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "get_payable_summary",
			"description": "Mendapatkan ringkasan hutang ke supplier.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties":  map[string]interface{}{},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "pay_receivable",
			"description": "Mencatat pembayaran piutang dari customer.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"customer_nama": map[string]interface{}{
						"type":        "string",
						"description": "Nama customer",
					},
					"jumlah": map[string]interface{}{
						"type":        "number",
						"description": "Jumlah yang dibayar dalam Rupiah",
					},
				},
				"required": []string{"customer_nama", "jumlah"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "pay_payable",
			"description": "Mencatat pembayaran hutang ke supplier.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"jumlah": map[string]interface{}{
						"type":        "number",
						"description": "Jumlah yang dibayar dalam Rupiah",
					},
				},
				"required": []string{"jumlah"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "toggle_wa_reminder",
			"description": "Mengaktifkan atau menonaktifkan pengiriman reminder WA ke customer.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"enabled": map[string]interface{}{
						"type":        "boolean",
						"description": "true untuk aktif, false untuk nonaktif",
					},
				},
				"required": []string{"enabled"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "list_sales_items",
			"description": "Menampilkan daftar semua item dalam database.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties":  map[string]interface{}{},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "list_sales_customers",
			"description": "Menampilkan daftar semua customer dalam database.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties":  map[string]interface{}{},
			},
		},
	},
}

// CreateSalesTransactionArgs represents arguments for create_sales_transaction.
type CreateSalesTransactionArgs struct {
	ItemNama     string `json:"item_nama"`
	Qty          int    `json:"qty"`
	CustomerNama string `json:"customer_nama"`
	Catatan      string `json:"catatan,omitempty"`
}

// AddSalesItemArgs represents arguments for add_sales_item.
type AddSalesItemArgs struct {
	Nama      string  `json:"nama"`
	HargaBeli float64 `json:"harga_beli"`
	Satuan    string  `json:"satuan,omitempty"`
}

// AddSalesCustomerArgs represents arguments for add_sales_customer.
type AddSalesCustomerArgs struct {
	Nama       string `json:"nama"`
	Alamat     string `json:"alamat"`
	Telepon    string `json:"telepon,omitempty"`
	JatuhTempo int    `json:"jatuh_tempo"`
	Payment    string `json:"payment"`
}

// SetCustomerPricingArgs represents arguments for set_customer_pricing.
type SetCustomerPricingArgs struct {
	CustomerNama string  `json:"customer_nama"`
	ItemNama     string  `json:"item_nama"`
	HargaJual    float64 `json:"harga_jual"`
}

// GetProfitReportArgs represents arguments for get_profit_report.
type GetProfitReportArgs struct {
	Period string `json:"period"`
}

// PayReceivableArgs represents arguments for pay_receivable.
type PayReceivableArgs struct {
	CustomerNama string  `json:"customer_nama"`
	Jumlah       float64 `json:"jumlah"`
}

// PayPayableArgs represents arguments for pay_payable.
type PayPayableArgs struct {
	Jumlah float64 `json:"jumlah"`
}

// ToggleWAReminderArgs represents arguments for toggle_wa_reminder.
type ToggleWAReminderArgs struct {
	Enabled bool `json:"enabled"`
}
