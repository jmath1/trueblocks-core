/*-------------------------------------------------------------------------
 * This source code is confidential proprietary information which is
 * copyright (c) 2018, 2019 TrueBlocks, LLC (http://trueblocks.io)
 * All Rights Reserved
 *------------------------------------------------------------------------*/
#include "options.h"

//-----------------------------------------------------------------------
bool COptions::exportCounts(void) {
    ENTER("exportCounts");

    bool isJson =
        (expContext().exportFmt == JSON1 || expContext().exportFmt == API1 || expContext().exportFmt == NONE1);

    bool first = true;
    for (auto aCount : counts) {
        if (isJson && !first)
            cout << ", ";
        cout << aCount.Format() << endl;
        first = false;
    }

    EXIT_NOMSG(true);
}