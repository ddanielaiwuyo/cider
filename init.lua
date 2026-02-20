-- Keybinds	
vim.g.mapleader = " "

vim.opt.signcolumn = "yes"
vim.opt.number = true
vim.opt.relativenumber = true
vim.opt.wrap = false
vim.opt.tabstop = 2
vim.opt.shiftwidth = 2
vim.opt.swapfile = false
vim.opt.winborder = "rounded"
vim.opt.smartindent = true
-- popup menu
vim.opt.pumheight = 10
vim.opt.completeopt = 'menu,menuone,noselect'
vim.opt.pumblend = 0

-- terminal config

-- escapae from terminal to normal mode
vim.keymap.set('t', '<Esc>', [[<C-\><C-n>]], { noremap = true })
vim.keymap.set('n', '<Esc>', ":nohl<CR>")

-- terminal mode
vim.keymap.set({ "n" }, "<leader>t", ":term<CR><C-wl>")

-- open explorer with space + e
vim.keymap.set('n', '<leader>e', ':Oil --float<CR>')
-- vim.keymap.set('n', '<leader>f', ':Pick files<CR>')
vim.keymap.set('n', '<leader>w', ':w<CR>')
vim.keymap.set('n', '<leader>ww', ':wa<CR>')

vim.keymap.set('n', '<Tab>', ':bn<CR>')
vim.keymap.set('n', '<leader>qq', ':qa!<CR>')
vim.keymap.set('n', '<leader>x', ':bdelete<CR>')


-- vim.keymap.set('n', '<leader>s', ':Lazy sync<CR>')
vim.keymap.set('n', '<leader>fm', vim.lsp.buf.format)
vim.keymap.set('n', '<leader>so', ':source ~/.config/nvim/init.lua<CR>')

-- ==================== for automcomplete navigation ==============
vim.keymap.set('i', '<Tab>', function()
	return vim.fn.pumvisible() == 1 and '<C-n>' or '<Tab>'
end, { expr = true })

vim.keymap.set('i', '<S-Tab>', function()
	return vim.fn.pumvisible() == 1 and '<C-p>' or '<S-Tab>'
end, { expr = true })

vim.keymap.set('i', '<CR>', function()
	return vim.fn.pumvisible() == 1 and '<C-y>' or '<CR>'
end, { expr = true })
--- ===============================================================


-- Using lazy.nvim
local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not (vim.uv or vim.loop).fs_stat(lazypath) then
	local lazyrepo = "https://github.com/folke/lazy.nvim.git"
	local out = vim.fn.system({ "git", "clone", "--filter=blob:none", "--branch=stable", lazyrepo, lazypath })
	if vim.v.shell_error ~= 0 then
		vim.api.nvim_echo({
			{ "[X] skilled issued hard, could not clone lazynvim\n", "ErrorMsg" },
			{ out,                                                   "WarningMsg" },
			{ "\n [i] get outta here" },
		}, true, {})
		vim.fn.getchar()
		os.exit(1)
	end
end
vim.opt.rtp:prepend(lazypath)


vim.o.cmdheight = 0

-- ======= plugins =========
require("lazy").setup({
	{
		"https://github.com/vague2k/vague.nvim",
		name = "vague",
		lazy = false,
		priority = 1000,
	},

	{
		"nvim-lualine/lualine.nvim",
		event = "VeryLazy",
		config = function()
			require("lualine").setup {
				options = {
					theme = "auto",
					section_seperators = '',
					component_seperators = '',
				},
				sections = {
					lualine_a = { "mode" },
					lualine_b = { "branch" },
					lualine_c = { { "filename", file_status = true } },
					-- lualine_c = { { "filename", file_status = true, path = 1 } },
					lualine_x = { "encoding", "filetype" },
					lualine_y = { "lsp_status" },
					-- lualine_z = { "time %{strftime('%H:%M')}" }, -- clock on the right
					lualine_z = {
						{ function()
							return ' ' .. os.date('%H:%M')
						end,
							color = { fg = '#4d7c86', bg = 'none', gui = 'bold' }
						},

						-- color = { bg = "none" },
					}, -- clock on the right

				},

				inactive_sections = {
					lualine_a = {},
					lualine_b = {},
					lualine_c = { "filename" },
					lualine_x = { "encoding", "filetype" },
					lualine_y = {},
					lualine_z = {},
				},
				tabline = {},
				extensions = {},
			}
		end
	},
	{
		'hrsh7th/nvim-cmp',
		dependencies = {
			'hrsh7th/cmp-nvim-lsp',
			'L3MON4D3/LuaSnip',
			'saadparwaiz1/cmp_luasnip',
		}
	},
	{
		"NTBBloodbath/doom-one.nvim",
		name = "doomchad",
	},
	{
		"navarasu/onedark.nvim",
		config = function()
			require("onedark").setup({
				style = 'darker',
			})
		end
	},
	{
		"https://github.com/neovim/nvim-lspconfig",
	},
	-- { "https://github.com/echasnovski/mini.pick" },
	{
		"nvim-telescope/telescope.nvim",
		lazy = false,
		priority = 1000,
		tag = "0.1.8",
		dependencies = { "nvim-lua/plenary.nvim" },
		-- config = function()
		-- 	vim.keymap.set("n", "<leader>ss", require("telescope.builtin").current_buffer_fuzzy_find
		-- 	)
		-- end
	},

	{ "https://github.com/stevearc/oil.nvim" },

	{
		"https://github.com/mason-org/mason.nvim",
		opts = {}
	}
})

-- --
-- -- Centered on screen
-- local win_config = function()
-- 	local height = math.floor(0.618 * vim.o.lines)
-- 	local width = math.floor(0.618 * vim.o.columns)
-- 	return {
-- 		anchor = 'NW',
-- 		height = height,
-- 		width = width,
-- 		row = math.floor(0.5 * (vim.o.lines - height)),
-- 		col = math.floor(0.5 * (vim.o.columns - width)),
-- 	}
-- end

-- require("mini.pick").setup({
-- 	mappings = {
-- 		move_up = "k",
-- 		move_down = "j",
-- 	},
--
-- 	{ window = { config = { border = 'double' } } },
-- 	{ window = { config = win_config } }
--
-- })
--


-- ====== telescope-config ======
require("telescope").setup({
	defaults = {
		-- file_ignore_patterns = { "node_modules", "%.git", "dist", "ven", "env", "__pycache_" },
		file_ignore_patterns = { "node_modules", "%.git", "ven", "env", "__pycache_", "target" },
		layout_strategy = "vertical",
		layout_config = { width = 0.7 }
		-- layout_config = {
		-- 	preview_width = 0.6,
		-- },
	}
})

local builtin = require("telescope.builtin")
vim.keymap.set("n", "<leader>ff", builtin.find_files, { desc = "Telescope find files" })
vim.keymap.set("n", "<leader>fb", builtin.buffers, { desc = "Run through buffers in telelscope" })
vim.keymap.set("n", "<leader>ss", builtin.live_grep, { desc = "Acts as live grep in cwd" })

-- set up quick fix for all errors in your shitty code
vim.keymap.set("n", "<leader>qf", builtin.diagnostics, { desc = "shows all dianogstics errors" })


-- ==== oil-config ====
require("oil").setup({
	default_file_explorer = true,
	buf_options = {
		buflisted = false,
		bufhidden = "hide",
	},
	float = {
		max_width = 0.6,
		max_height = 0.6,
	},
	columns = { " icon" },
	sort = {
		sorter = "case_sensitive",
		directories_first = true,
	},

	keymaps = {
		["q"] = { "actions.close", mode = 'n' },
	},

	win_options = {
		number = false,
		relativenumber = false,
		spell = false,
		signcolumn = "yes:2"
	},

	view_options = {
		show_hidden = true,
	}


})


vim.lsp.enable({ "lua_ls", "gopls", "clangd", "ts_ls", "jdtls", "pyright", "svelteserver", "html", "vtsls" })
-- Tells lua that vim is a global variable
vim.lsp.config("lua_ls", {
	settings = {
		Lua = {
			workspace = {
				library = vim.api.nvim_get_runtime_file("", true)
			}
		}
	}
})

vim.lsp.config("svelteserver", {})

-- Autocompletion
vim.api.nvim_create_autocmd('LspAttach', {
	callback = function(ev)
		local client = vim.lsp.get_client_by_id(ev.data.client_id)
		if client:supports_method('textDocument/completion') then
			vim.lsp.completion.enable(true, client.id, ev.buf, { autotrigger = true })
		end

		vim.keymap.set('n', 'K', function()
			vim.lsp.buf.hover({ border = 'rounded' })
		end, { buffer = ev.buf }
		)

		vim.keymap.set('n', '<C-k>', function()
			vim.lsp.buf.signature_help({ border = 'rounded' })
		end, { buffer = ev.buf })
	end,
})

vim.cmd("set  completeopt+=noselect")


vim.api.nvim_create_autocmd("TextYankPost", {
	callback = function()
		vim.highlight.on_yank()
	end,
})



-- ==== lsp dianogstics ======
vim.diagnostic.config({
	virtual_text = {
		spacing = 6,
		prefix = '*',
	},
	signs = true,
	float = {
		border = "rounded",
	},
	update_in_insert = false,
	severity_sort = true,
})

vim.cmd.colorscheme("vague")



-- Setup LuaSnip
local luasnip = require('luasnip')

-- Setup nvim-cmp
local cmp = require('cmp')
cmp.setup({
	snippet = {
		expand = function(args)
			luasnip.lsp_expand(args.body)
		end,
	},
	mapping = {
		['<Tab>'] = cmp.mapping(function(fallback)
			if cmp.visible() then
				cmp.select_next_item()
			elseif luasnip.expand_or_jumpable() then
				luasnip.expand_or_jump()
			else
				fallback()
			end
		end, { 'i', 's' }),
		['<S-Tab>'] = cmp.mapping(function(fallback)
			if cmp.visible() then
				cmp.select_prev_item()
			elseif luasnip.jumpable(-1) then
				luasnip.jump(-1)
			else
				fallback()
			end
		end, { 'i', 's' }),
		['<CR>'] = cmp.mapping.confirm({ select = true }),
	},
	sources = {
		{ name = 'nvim_lsp' },
		{ name = 'luasnip' },
	},
})

-- vim.o.statusline = "%f %=%{strftime('%H:%M')} 🫀"
