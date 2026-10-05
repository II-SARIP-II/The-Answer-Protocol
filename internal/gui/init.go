package gui

import (
    "tap/internal/gui/components"
    "tap/internal/gui/components/panels"
)

func InitPanels(winWidth, winHeight int) (*components.Panel, *components.Panel, *components.Panel, *components.Panel, *components.Panel, *components.Panel, *components.Panel, *components.Panel) {
    w := float32(winWidth)
    h := float32(winHeight)

    topH := h * 0.08
    centerH := h * 0.55
    bottomH := h - topH - centerH

    colLeftW := w * 0.32
    colMidW := w * 0.48
    colRightW := w - colLeftW - colMidW

    t_panel := panels.CreateTopPanel(w, topH)
    cl_panel := panels.CreatePlayerPanel(topH, centerH, colLeftW)
    cm_panel := panels.CreateMapPanel(topH, centerH, colLeftW, colMidW)
    cr_panel := panels.CreateActionsPanel(topH, centerH, colLeftW, colMidW, colRightW)
    bl_panel := panels.CreateStatsPanel(topH, centerH, bottomH, colLeftW)
    bcl_panel := panels.CreateStuffPanel(topH, centerH, bottomH, colLeftW)
    bcr_panel := panels.CreateChatPanel(topH, centerH, bottomH, colLeftW, colMidW)
    br_panel := panels.CreateInteractionsPanel(topH, centerH, bottomH, colLeftW, colMidW, colRightW)

    return t_panel, cl_panel, cm_panel, cr_panel, bl_panel, bcl_panel, bcr_panel, br_panel
}
