// +build f469xx
package display

import "github.com/ianmcmahon/emgo/egpath/src/stm32/o/f469xx/dsi"

/**
 * @brief  Initializes the DSI LCD.
 * @retval LCD state
 */
func BSP_LCD_Init() uint8 {
	return BSP_LCD_InitEx(LCD_ORIENTATION_LANDSCAPE)
}

/**
 * @brief  Initializes the DSI LCD.
 * The ititialization is done as below:
 *     - DSI PLL ititialization
 *     - DSI ititialization
 *     - LTDC ititialization
 *     - OTM8009A LCD Display IC Driver ititialization
 * @retval LCD state
 */
func BSP_LCD_InitEx(LCD_OrientationTypeDef orientation) uint8 {
	dsi := dsi.DSI
	var PhyTimings DSI_PHY_TimerTypeDef

	// static RCC_PeriphCLKInitTypeDef  PeriphClkInitStruct
	var PeriphClkInitStruct RCC_PeriphCLKInitTypeDef

	var LcdClock uint32_t = 27429 /*!< LcdClk = 27429 kHz */

	var laneByteClk_kHz = 0
	var VSA uint32  /*!< Vertical start active time in units of lines */
	var VBP uint32  /*!< Vertical Back Porch time in units of lines */
	var VFP uint32  /*!< Vertical Front Porch time in units of lines */
	var VACT uint32 /*!< Vertical Active time in units of lines = imageSize Y in pixels to display */
	var HSA uint32  /*!< Horizontal start active time in units of lcdClk */
	var HBP uint32  /*!< Horizontal Back Porch time in units of lcdClk */
	var HFP uint32  /*!< Horizontal Front Porch time in units of lcdClk */
	var HACT uint32 /*!< Horizontal Active time in units of lcdClk = imageSize X in pixels to display */

	/* Toggle Hardware Reset of the DSI LCD using
	 * its XRES signal (active low) */
	BSP_LCD_Reset()

	/* Call first MSP Initialize only in case of first initialization
	 * This will set IP blocks LTDC, DSI and DMA2D
	 * - out of reset
	 * - clocked
	 * - NVIC IRQ related to IP blocks enabled
	 */
	BSP_LCD_MspInit()

	/*************************DSI Initialization***********************************/

	/* Base address of DSI Host/Wrapper registers to be set before calling De-Init */
	hdsi_eval.Instance = DSI

	HAL_DSI_DeInit(&hdsi_eval)

	// TODO: check for revA, uses different numbers
	dsiPllInit.PLLNDIV = 125
	dsiPllInit.PLLIDF = DSI_PLL_IN_DIV2
	dsiPllInit.PLLODF = DSI_PLL_OUT_DIV1
	laneByteClk_kHz = 62500 /* 500 MHz / 8 = 62.5 MHz = 62500 kHz */

	/* Set number of Lanes */
	hdsi_eval.Init.NumberOfLanes = DSI_TWO_DATA_LANES

	/* TXEscapeCkdiv = f(LaneByteClk)/15.62 = 4 */
	hdsi_eval.Init.TXEscapeCkdiv = laneByteClk_kHz / 15620

	HAL_DSI_Init(&(hdsi_eval), &(dsiPllInit))

	/* Timing parameters for all Video modes
	 * Set Timing parameters of LTDC depending on its chosen orientation
	 */
	if orientation == LCD_ORIENTATION_PORTRAIT {
		lcd_x_size = OTM8009A_480X800_WIDTH  /* 480 */
		lcd_y_size = OTM8009A_480X800_HEIGHT /* 800 */
	} else {
		/* lcd_orientation == LCD_ORIENTATION_LANDSCAPE */
		lcd_x_size = OTM8009A_800X480_WIDTH /* 800 */
		lcd_y_size = OTM8009A_800X480_HEIGH /* 480 */
	}

	HACT = lcd_x_size
	VACT = lcd_y_size

	/* The following values are same for portrait and landscape orientations */
	VSA = OTM8009A_480X800_VSYNC
	VBP = OTM8009A_480X800_VBP
	VFP = OTM8009A_480X800_VFP
	HSA = OTM8009A_480X800_HSYNC
	HBP = OTM8009A_480X800_HBP
	HFP = OTM8009A_480X800_HFP

	hdsivideo_handle.VirtualChannelID = LCD_OTM8009A_ID
	hdsivideo_handle.ColorCoding = LCD_DSI_PIXEL_DATA_FMT_RBG888
	hdsivideo_handle.VSPolarity = DSI_VSYNC_ACTIVE_HIGH
	hdsivideo_handle.HSPolarity = DSI_HSYNC_ACTIVE_HIGH
	hdsivideo_handle.DEPolarity = DSI_DATA_ENABLE_ACTIVE_HIGH
	hdsivideo_handle.Mode = DSI_VID_MODE_BURST /* Mode Video burst ie : one LgP per line */
	hdsivideo_handle.NullPacketSize = 0xFFF
	hdsivideo_handle.NumberOfChunks = 0
	hdsivideo_handle.PacketSize = HACT /* Value depending on display orientation choice portrait/landscape */
	hdsivideo_handle.HorizontalSyncActive = (HSA * laneByteClk_kHz) / LcdClock
	hdsivideo_handle.HorizontalBackPorch = (HBP * laneByteClk_kHz) / LcdClock
	hdsivideo_handle.HorizontalLine = ((HACT + HSA + HBP + HFP) * laneByteClk_kHz) / LcdClock /* Value depending on display orientation choice portrait/landscape */
	hdsivideo_handle.VerticalSyncActive = VSA
	hdsivideo_handle.VerticalBackPorch = VBP
	hdsivideo_handle.VerticalFrontPorch = VFP
	hdsivideo_handle.VerticalActive = VACT /* Value depending on display orientation choice portrait/landscape */

	/* Enable or disable sending LP command while streaming is active in video mode */
	hdsivideo_handle.LPCommandEnable = DSI_LP_COMMAND_ENABLE /* Enable sending commands in mode LP (Low Power) */

	/* Largest packet size possible to transmit in LP mode in VSA, VBP, VFP regions */
	/* Only useful when sending LP packets is allowed while streaming is active in video mode */
	hdsivideo_handle.LPLargestPacketSize = 16

	/* Largest packet size possible to transmit in LP mode in HFP region during VACT period */
	/* Only useful when sending LP packets is allowed while streaming is active in video mode */
	hdsivideo_handle.LPVACTLargestPacketSize = 0

	/* Specify for each region of the video frame, if the transmission of command in LP mode is allowed in this region */
	/* while streaming is active in video mode                                                                         */
	hdsivideo_handle.LPHorizontalFrontPorchEnable = DSI_LP_HFP_ENABLE /* Allow sending LP commands during HFP period */
	hdsivideo_handle.LPHorizontalBackPorchEnable = DSI_LP_HBP_ENABLE  /* Allow sending LP commands during HBP period */
	hdsivideo_handle.LPVerticalActiveEnable = DSI_LP_VACT_ENABLE      /* Allow sending LP commands during VACT period */
	hdsivideo_handle.LPVerticalFrontPorchEnable = DSI_LP_VFP_ENABLE   /* Allow sending LP commands during VFP period */
	hdsivideo_handle.LPVerticalBackPorchEnable = DSI_LP_VBP_ENABLE    /* Allow sending LP commands during VBP period */
	hdsivideo_handle.LPVerticalSyncActiveEnable = DSI_LP_VSYNC_ENABLE /* Allow sending LP commands during VSync = VSA period */

	/* Configure DSI Video mode timings with settings set above */
	HAL_DSI_ConfigVideoMode(&(hdsi_eval), &(hdsivideo_handle))

	/* Configure DSI PHY HS2LP and LP2HS timings */
	PhyTimings.ClockLaneHS2LPTime = 35
	PhyTimings.ClockLaneLP2HSTime = 35
	PhyTimings.DataLaneHS2LPTime = 35
	PhyTimings.DataLaneLP2HSTime = 35
	PhyTimings.DataLaneMaxReadTime = 0
	PhyTimings.StopWaitTime = 10
	HAL_DSI_ConfigPhyTimer(&hdsi_eval, &PhyTimings)

	/*************************End DSI Initialization*******************************/

	/************************LTDC Initialization***********************************/

	/* Timing Configuration */
	hltdc_eval.Init.HorizontalSync = (HSA - 1)
	hltdc_eval.Init.AccumulatedHBP = (HSA + HBP - 1)
	hltdc_eval.Init.AccumulatedActiveW = (lcd_x_size + HSA + HBP - 1)
	hltdc_eval.Init.TotalWidth = (lcd_x_size + HSA + HBP + HFP - 1)

	/* Initialize the LCD pixel width and pixel height */
	hltdc_eval.LayerCfg.ImageWidth = lcd_x_size
	hltdc_eval.LayerCfg.ImageHeight = lcd_y_size

	/* LCD clock configuration */
	/* PLLSAI_VCO Input = HSE_VALUE/PLL_M = 1 Mhz */
	/* PLLSAI_VCO Output = PLLSAI_VCO Input * PLLSAIN = 384 Mhz */
	/* PLLLCDCLK = PLLSAI_VCO Output/PLLSAIR = 384 MHz / 7 = 54.857 MHz */
	/* LTDC clock frequency = PLLLCDCLK / LTDC_PLLSAI_DIVR_2 = 54.857 MHz / 2 = 27.429 MHz */
	PeriphClkInitStruct.PeriphClockSelection = RCC_PERIPHCLK_LTDC
	PeriphClkInitStruct.PLLSAI.PLLSAIN = 384
	PeriphClkInitStruct.PLLSAI.PLLSAIR = 7
	PeriphClkInitStruct.PLLSAIDivR = RCC_PLLSAIDIVR_2
	HAL_RCCEx_PeriphCLKConfig(&PeriphClkInitStruct)

	/* Background value */
	hltdc_eval.Init.Backcolor.Blue = 0
	hltdc_eval.Init.Backcolor.Green = 0
	hltdc_eval.Init.Backcolor.Red = 0
	hltdc_eval.Init.PCPolarity = LTDC_PCPOLARITY_IPC
	hltdc_eval.Instance = LTDC

	/* Get LTDC Configuration from DSI Configuration */
	HAL_LTDCEx_StructInitFromVideoConfig(&(hltdc_eval), &(hdsivideo_handle))

	/* Initialize the LTDC */
	HAL_LTDC_Init(&hltdc_eval)

	/* Enable the DSI host and wrapper after the LTDC initialization
	   To avoid any synchronization issue, the DSI shall be started after enabling the LTDC */
	HAL_DSI_Start(&(hdsi_eval))

	// TODO: make this enableable
	/* Initialize the SDRAM */
	BSP_SDRAM_Init()

	/* Initialize the font */
	BSP_LCD_SetFont(&LCD_DEFAULT_FONT)

	/************************End LTDC Initialization*******************************/

	/***********************OTM8009A Initialization********************************/

	/* Initialize the OTM8009A LCD Display IC Driver (KoD LCD IC Driver)
	 *  depending on configuration set in 'hdsivideo_handle'.
	 */
	OTM8009A_Init(OTM8009A_FORMAT_RGB888, orientation)

	/***********************End OTM8009A Initialization****************************/

	return LCD_OK
}
